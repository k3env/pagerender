package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gofiber/fiber/v3"
)

type RenderParams struct {
	Url      string `json:"url"`
	FileName string `json:"filename"`
}

func RenderHandler() fiber.Handler {
	var margin float64 = 0.05
	var ppi float64 = 96

	return func(ctx fiber.Ctx) error {
		var params RenderParams
		err := ctx.Bind().JSON(&params)

		if err != nil {
			return ctx.Status(http.StatusBadRequest).SendString(err.Error())
		}

		browserURL := launcher.New().
			Headless(true).
			Set(flags.NoSandbox, "").
			MustLaunch()

		browser := rod.New().ControlURL(browserURL).MustConnect()
		defer browser.MustClose()

		page := browser.MustPage(params.Url)
		defer page.MustClose()

		page.MustWaitLoad()
		page.MustWaitStable()

		time.Sleep(3 * time.Second)

		container, err := page.Element("#render-container")
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString("Render container not found")
		}

		box, err := container.Shape()
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString("Cant get shape of container")
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
			log.Println("PDF error:", err)
			return ctx.Status(500).SendString("Failed to generate PDF")
		}

		bytes, err := io.ReadAll(pdf)

		// Отдаём PDF
		ctx.Set("Content-Type", "application/pdf")
		ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", params.FileName))
		return ctx.Send(bytes)
	}
}
