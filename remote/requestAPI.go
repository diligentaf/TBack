package remote

import (
	"TBack/conf"

	"time"

	"github.com/juju/errors"
	//resty "gopkg.in/resty.v1"
	resty "github.com/go-resty/resty"
)

// GetAPI ...
func GetAPI(trID, url string) ([]byte, error) {
	mlog.Infow("GetAPI", "trID", trID, "url", url)

	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetBasicAuth(conf.TBack.GetString("remote_user"), conf.TBack.GetString("remote_pass"))

	resp, err := client.R().
		SetHeader("X-Request-ID", trID).
		Get(url)
	if err != nil {
		return nil, errors.Annotate(err, "remote GetAPI")
	}
	if status := resp.StatusCode(); status != 200 {
		return nil, errors.BadRequestf("remote GetAPI: Response status code isn't 200 statusCode[%d]", status)
	}

	if resp.Time() > 10*time.Second {
		mlog.Warnw("remote GetAPI Long Response Time", "trID", trID, "elapsed", resp.Time().String(), "url", url)
	}

	return resp.Body(), nil
}

// PostAPI ...
func PostAPI(trID, url string, req []byte) ([]byte, error) {
	mlog.Infow("PostAPI request", "trID", trID, "reqBodyLen", len(req), "url", url)

	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetBasicAuth(conf.TBack.GetString("remote_user"), conf.TBack.GetString("remote_pass"))

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Request-ID", trID).
		SetBody(req).
		Post(url)
	if err != nil {
		return nil, errors.Annotate(err, "remote PostAPI")
	}
	if status := resp.StatusCode(); status != 200 {
		return nil, errors.BadRequestf("remote PostAPI: Response status code isn't 200 statusCode[%d]", status)
	}

	if resp.Time() > 3*time.Second {
		mlog.Warnw("remote PostAPI Long Response Time", "trID", trID, "elapsed", resp.Time().String(), "reqBodyLen", len(req), "url", url)
	}

	return resp.Body(), nil
}
