package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gofiber/fiber/v3"
)

var (
	margin float64 = 0.05
	ppi    float64 = 96

	ErrNoShape     = errors.New("no shape found")
	ErrNoContainer = errors.New("no container found")
)

type RenderParams struct {
	Url      string `json:"url"`
	FileName string `json:"filename"`
}

func RenderHandler() fiber.Handler {

	return func(ctx fiber.Ctx) error {
		start := time.Now()
		var params RenderParams
		err := ctx.Bind().JSON(&params)

		if err != nil {
			return ctx.Status(http.StatusBadRequest).SendString(err.Error())
		}

		bytes, err := renderToBytes(params)
		if err != nil {
			pdfRenderFailures.Inc()
			if errors.Is(err, ErrNoShape) || errors.Is(err, ErrNoContainer) {
				return ctx.Status(http.StatusBadRequest).SendString(err.Error())
			}
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		// Отдаём PDF
		ctx.Set("Content-Type", "application/pdf")
		ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", params.FileName))
		pdfRenderDuration.Observe(time.Since(start).Seconds())
		return ctx.Send(bytes)
	}
}

func renderToBytes(params RenderParams) ([]byte, error) {
	browserURL, err := launcher.New().
		Headless(true).
		Set(flags.NoSandbox, "").
		Launch()

	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(browserURL)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{URL: params.Url})
	if err != nil {
		return nil, err
	}
	defer page.MustClose()

	err = page.WaitLoad()
	if err != nil {
		return nil, err
	}
	err = page.WaitStable(3 * time.Second)
	if err != nil {
		return nil, err
	}

	time.Sleep(3 * time.Second)

	container, err := page.Element("#render-container")
	if err != nil {
		return nil, ErrNoContainer
	}

	box, err := container.Shape()
	if err != nil {
		return nil, ErrNoShape
	}

	wi := box.Box().Width / ppi
	hi := box.Box().Height / ppi

	// Генерация PDF
	pdf, err := page.PDF(&proto.PagePrintToPDF{
		PrintBackground: true,
		PaperWidth:      &wi,
		PaperHeight:     &hi,
		MarginTop:       &margin,
		MarginBottom:    &margin,
		MarginLeft:      &margin,
		MarginRight:     &margin,
		Landscape:       false,
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(pdf)
}
