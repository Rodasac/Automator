package browser_automator

import (
	"automator-go/entities/models"
	"automator-go/usecases/task"
	"fmt"
	"github.com/go-rod/rod"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type RodAutomator struct {
	browser  *rod.Browser
	pagePool rod.PagePool
	logger   *zap.Logger
}

func NewRodAutomator(browser *rod.Browser, pagePool rod.PagePool, logger *zap.Logger) *RodAutomator {
	return &RodAutomator{browser: browser, pagePool: pagePool, logger: logger}
}

func (at *RodAutomator) createPage() *rod.Page {
	return at.browser.MustPage()
}

func (at *RodAutomator) Run(taskToRun *models.Task) (*[]task.RawMedia, error) {
	at.logger.Debug("Getting page from pool")
	page := at.pagePool.Get(at.createPage)
	defer at.pagePool.Put(page)

	err := page.Navigate(taskToRun.Url)
	if err != nil {
		return nil, fmt.Errorf("error navigating to url: %w", err)
	}
	at.logger.Debug("Page initialized and navigated to url")

	pageTimeOutEnv := os.Getenv("BROWSER_PAGE_TIMEOUT_BY_TASK")
	if strings.TrimSpace(pageTimeOutEnv) == "" {
		pageTimeOutEnv = "15s"
	}
	pageTimeout, err := time.ParseDuration(pageTimeOutEnv)
	if err != nil {
		return nil, fmt.Errorf("error parsing page timeout: %w", err)
	}
	page = page.Timeout(pageTimeout)
	at.logger.Debug("Set page timeout")

	err = page.WaitStable(800*time.Millisecond, 1)
	if err != nil {
		return nil, fmt.Errorf("error waiting for page to be stable: %w", err)
	}
	at.logger.Debug("Page is stable and loaded")

	rawMedias := make([]task.RawMedia, 0)

	for _, action := range taskToRun.Actions {
		switch action.Type {
		case models.Navigate:
			at.logger.Debug("Navigating to url", zap.String("selector", action.Value))
			err = navigate(page, action)
			if err != nil {
				return nil, fmt.Errorf("error navigating xpath: %w", err)
			}
			at.logger.Debug("Navigated to url", zap.String("selector", action.Value))
		case models.Click:
			at.logger.Debug("Clicking on element", zap.String("selector", action.Value))
			err = click(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clicking xpath: %w", err)
			}
			at.logger.Debug("Clicked on element", zap.String("selector", action.Value))
		case models.ScrollDown:
			at.logger.Debug("Scrolling down")
			err = scrollDown(page, action)
			if err != nil {
				return nil, fmt.Errorf("error scrolling down: %w", err)
			}
			at.logger.Debug("Scrolled down")
		case models.Capture:
			at.logger.Debug("Capturing element", zap.String("selector", action.Value))
			rawMedia, err := capture(page, action)
			if err != nil {
				return nil, fmt.Errorf("error capturing element: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
			at.logger.Debug("Captured element", zap.String("selector", action.Value))
		case models.WaitSeconds:
			at.logger.Debug("Waiting seconds", zap.String("seconds", action.Value))
			err = waitSeconds(action)
			if err != nil {
				return nil, fmt.Errorf("error waiting seconds: %w", err)
			}
			at.logger.Debug("Waited seconds", zap.String("seconds", action.Value))
		case models.WriteInput:
			at.logger.Debug("Writing input", zap.String("input", action.Value))
			err = writeInput(page, action)
			if err != nil {
				return nil, fmt.Errorf("error writing input: %w", err)
			}
			at.logger.Debug("Wrote input", zap.String("input", action.Value))
		case models.ClearInput:
			at.logger.Debug("Clearing input")
			err = clearInput(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clearing input: %w", err)
			}
			at.logger.Debug("Cleared input")
		case models.SelectOptions:
			at.logger.Debug("Selecting options", zap.String("options", action.Value))
			err = selectOptions(page, action)
			if err != nil {
				return nil, fmt.Errorf("error selecting options: %w", err)
			}
			at.logger.Debug("Selected options", zap.String("options", action.Value))
		case models.WriteTime:
			at.logger.Debug("Writing time on input", zap.String("input", action.Value))
			err = writeTime(page, action)
			if err != nil {
				return nil, fmt.Errorf("error writing time on input: %w", err)
			}
			at.logger.Debug("Wrote time on input", zap.String("input", action.Value))
		case models.DownloadResource:
			at.logger.Debug("Downloading resource", zap.String("selector", action.Value))
			rawMedia, err := downloadResource(page, action)
			if err != nil {
				return nil, fmt.Errorf("error downloading resource: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
			at.logger.Debug("Downloaded resource", zap.String("selector", action.Value))
		default:
			return nil, fmt.Errorf("unknown action type: %s", action.Type.String())
		}
	}

	return &rawMedias, nil
}
