package main

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestPriceDataSma8(t *testing.T) {
	prices := [8]float64{7, 3, 5, 1, 9, 3, 4, 5}
	pdList := PriceDataList{}

	for _, price := range prices {
		pdList = append(pdList, PriceData{price})
	}

	res := pdList.sma(8, 0)
	expected := 4.625

	assert.Equal(t, res, expected)
}

func TestPriceDataSm8Offset(t *testing.T) {
	prices := [10]float64{7, 3, 5, 1, 9, 3, 4, 5, 8, 9} //We add two extra values which should be ignored in our calculation
	pdList := PriceDataList{}

	for _, price := range prices {
		pdList = append(pdList, PriceData{price})
	}

	res := pdList.sma(8, 2)
	expected := 4.625

	assert.Equal(t, res, expected)
}

func TestParsePeriod(t *testing.T) {
	res1m, _ := parsePeriod("1m")
	res5m, _ := parsePeriod("5m")
	res15m, _ := parsePeriod("15m")
	res30m, _ := parsePeriod("30m")
	res1h, _ := parsePeriod("1h")
	res2h, _ := parsePeriod("2h")
	res4h, _ := parsePeriod("4h")
	res1d, _ := parsePeriod("1d")
	_, err := parsePeriod("20d")

	assert.Equal(t, time.Minute, res1m)
	assert.Equal(t, time.Minute*5, res5m)
	assert.Equal(t, time.Minute*15, res15m)
	assert.Equal(t, time.Minute*30, res30m)
	assert.Equal(t, time.Hour, res1h)
	assert.Equal(t, time.Hour*2, res2h)
	assert.Equal(t, time.Hour*4, res4h)
	assert.Equal(t, time.Hour*24, res1d)
	assert.Equal(t, err, errors.New("Invalid period 20d"))
}

func TestFetchPriceData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponseBody := `
		[{"Open":35426.88,"High":35620.52,"Low":35180.72,"Close":35403.19,"BaseVolume":24.36630896,"QuoteVolume":0.0011324399691342,"OpenTime":"2021-05-18T21:00:00Z"},{"Open":35403.19,"High":35416.69,"Low":35021.21,"Close":35192.79,"BaseVolume":25.718721819999995,"QuoteVolume":0.0010189815988635,"OpenTime":"2021-05-18T21:30:00Z"},{"Open":35192.79,"High":35338.29,"Low":34976.53,"Close":35306.78,"BaseVolume":36.347537069999994,"QuoteVolume":0.0010619061358844999,"OpenTime":"2021-05-18T22:00:00Z"},{"Open":35306.78,"High":35425.85,"Low":35020.13,"Close":35052.38,"BaseVolume":25.370089730000004,"QuoteVolume":0.0011111079448168,"OpenTime":"2021-05-18T22:30:00Z"},{"Open":35052.38,"High":35096.96,"Low":34609.61,"Close":35096.96,"BaseVolume":48.86442366000001,"QuoteVolume":0.0021593336079350997,"OpenTime":"2021-05-18T23:00:00Z"},{"Open":35102.56,"High":35355.04,"Low":34716.8,"Close":35090.95,"BaseVolume":37.72534407999999,"QuoteVolume":0.0017416981124196,"OpenTime":"2021-05-18T23:30:00Z"},{"Open":35090.95,"High":35662.71,"Low":34809.79,"Close":35369.76,"BaseVolume":31.564842559999995,"QuoteVolume":0.0015733705989539998,"OpenTime":"2021-05-19T00:00:00Z"},{"Open":35367.65,"High":35371.06,"Low":34836.56,"Close":34885.95,"BaseVolume":26.36056495,"QuoteVolume":0.0012521697085931998,"OpenTime":"2021-05-19T00:30:00Z"},{"Open":34885.95,"High":34967.85,"Low":33875.95,"Close":34009.53,"BaseVolume":94.34789558000001,"QuoteVolume":0.004378596227738801,"OpenTime":"2021-05-19T01:00:00Z"},{"Open":33997.12,"High":34158.85,"Low":33185,"Close":33470.93,"BaseVolume":88.36032580999999,"QuoteVolume":0.0045237149070227,"OpenTime":"2021-05-19T01:30:00Z"},{"Open":33470.93,"High":33849.81,"Low":33294.23,"Close":33470.08,"BaseVolume":79.34479463999999,"QuoteVolume":0.0020865520408099,"OpenTime":"2021-05-19T02:00:00Z"},{"Open":33470.08,"High":33699.75,"Low":33030.56,"Close":33045.38,"BaseVolume":52.878230980000005,"QuoteVolume":0.0026875057208354,"OpenTime":"2021-05-19T02:30:00Z"},{"Open":33049.54,"High":33608.52,"Low":32858.23,"Close":33554.05,"BaseVolume":110.95060591000002,"QuoteVolume":0.0034325613010108003,"OpenTime":"2021-05-19T03:00:00Z"},{"Open":33554.05,"High":33557.96,"Low":33095.27,"Close":33220,"BaseVolume":33.11700903999999,"QuoteVolume":0.0017339097021018,"OpenTime":"2021-05-19T03:30:00Z"},{"Open":33220,"High":33433.24,"Low":32250,"Close":32311.11,"BaseVolume":187.57766748,"QuoteVolume":0.006746458146826799,"OpenTime":"2021-05-19T04:00:00Z"},{"Open":32311.11,"High":32798.22,"Low":31500,"Close":32059.71,"BaseVolume":291.11489882999996,"QuoteVolume":0.013990895300625302,"OpenTime":"2021-05-19T04:30:00Z"},{"Open":32065.39,"High":34966,"Low":32011.81,"Close":32494.61,"BaseVolume":290.99993481000007,"QuoteVolume":0.009468383941342499,"OpenTime":"2021-05-19T05:00:00Z"},{"Open":32494.61,"High":32497.32,"Low":31846.67,"Close":32189.79,"BaseVolume":149.12703732,"QuoteVolume":0.0071526185417292,"OpenTime":"2021-05-19T05:30:00Z"},{"Open":32012.76,"High":32580.6,"Low":32006.66,"Close":32506.89,"BaseVolume":167.39679049999998,"QuoteVolume":0.0051574763115506995,"OpenTime":"2021-05-19T06:00:00Z"},{"Open":32506.89,"High":32506.9,"Low":31929.58,"Close":32131.94,"BaseVolume":100.98459516999998,"QuoteVolume":0.004904070476965501,"OpenTime":"2021-05-19T06:30:00Z"},{"Open":32151.25,"High":32376.92,"Low":31473.58,"Close":32250,"BaseVolume":228.21320169999998,"QuoteVolume":0.008020366763716603,"OpenTime":"2021-05-19T07:00:00Z"},{"Open":32250,"High":33146.43,"Low":32140.31,"Close":33033.23,"BaseVolume":195.44355966999998,"QuoteVolume":0.008880698753919402,"OpenTime":"2021-05-19T07:30:00Z"},{"Open":33010.5,"High":33080,"Low":32603.33,"Close":32978.91,"BaseVolume":136.19600679,"QuoteVolume":0.0048661444907722,"OpenTime":"2021-05-19T08:00:00Z"},{"Open":32978.91,"High":33372.65,"Low":32740.79,"Close":33101.23,"BaseVolume":123.16877149,"QuoteVolume":0.005470036168446601,"OpenTime":"2021-05-19T08:30:00Z"},{"Open":33101.23,"High":33310.21,"Low":32780.55,"Close":32851.88,"BaseVolume":73.45633537999998,"QuoteVolume":0.0028880525914528996,"OpenTime":"2021-05-19T09:00:00Z"},{"Open":32851.88,"High":33084.93,"Low":32739.54,"Close":32997.75,"BaseVolume":67.09209892000001,"QuoteVolume":0.0030597544088610006,"OpenTime":"2021-05-19T09:30:00Z"},{"Open":32986.48,"High":33197.23,"Low":32157.05,"Close":32362.57,"BaseVolume":138.42558864999995,"QuoteVolume":0.0058010884380137,"OpenTime":"2021-05-19T10:00:00Z"},{"Open":32362.57,"High":32406.33,"Low":31765,"Close":32180.56,"BaseVolume":247.21953638,"QuoteVolume":0.007817722468304,"OpenTime":"2021-05-19T10:30:00Z"},{"Open":32149.58,"High":32344.92,"Low":30763,"Close":30818.57,"BaseVolume":248.47577595000004,"QuoteVolume":0.0116233866964768,"OpenTime":"2021-05-19T11:00:00Z"},{"Open":30818.57,"High":31400,"Low":28566,"Close":31342.4,"BaseVolume":481.59251555000003,"QuoteVolume":0.022728913454508292,"OpenTime":"2021-05-19T11:30:00Z"},{"Open":31280.49,"High":31820,"Low":30878.92,"Close":31055.2,"BaseVolume":212.08216445,"QuoteVolume":0.009912728835866398,"OpenTime":"2021-05-19T12:00:00Z"}]
	 `

	httpmock.RegisterResponder(
		"GET",
		`=~^http://cryptohopper-ticker-frontend\.us-east-1\.elasticbeanstalk\.com/v1/coinbasepro/candles\?pair=BTC-EUR&start=\d+&end=\d+&period=30m\z`,
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200, mockResponseBody)
			res.Header.Set("Content-Type", "application/json")

			return res, nil
		},
	)

	pdl, _ := fetchPriceData("coinbasepro", "BTC-EUR", "30m")

	assert.Equal(t, 31, len(pdl))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

}

type SMAMockCall struct {
	n         int
	offset    int
	returnVal float64
}

func mockSMACalls(ctrl *gomock.Controller, mockCalls []SMAMockCall) SMA {
	n := NewMockSMA(ctrl)

	for _, mc := range mockCalls {
		n.
			EXPECT().
			sma(mc.n, mc.offset).
			Return(mc.returnVal)
	}

	return n
}

func TestGenerateSignalBuy(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockCalls := [3]SMAMockCall{
		{8, 0, 500},  //Current SMA(8)
		{8, 1, 400},  //Previous SMA(8)
		{55, 0, 450}, // SMA(55)
	}

	n := mockSMACalls(ctrl, mockCalls[:])

	res := generateSignal(n)

	assert.Equal(t, "BUY", res)
}

func TestGenerateSignalSell(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockCalls := [3]SMAMockCall{
		{8, 0, 400},  //Current SMA(8)
		{8, 1, 450},  //Previous SMA(8)
		{55, 0, 450}, // SMA(55)
	}

	n := mockSMACalls(ctrl, mockCalls[:])

	res := generateSignal(n)

	assert.Equal(t, "SELL", res)
}

func TestGenerateSignalNeutral(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockCalls := [3]SMAMockCall{
		{8, 0, 400},  //Current SMA(8)
		{8, 1, 400},  //Previous SMA(8)
		{55, 0, 400}, // SMA(55)
	}

	n := mockSMACalls(ctrl, mockCalls[:])

	res := generateSignal(n)

	assert.Equal(t, "NEUTRAL", res)
}

//func TestSignalBuy(t *testing.T) {
//	current_sma8 := 50
//	previous_sma8 := 49
//	sma55 := 40
//
//	url := "/signal?period=5m&exchange=coinbasepro&pair=BTC-EUR"
//
//	req := httptest.NewRequest(http.MethodGet, url, nil)
//	w := httptest.NewRecorder()
//	signal(w, req)
//
//	res := w.Result()
//	defer res.Body.Close()
//
//	data, err := ioutil.ReadAll(res.Body)
//
//	//mock the cryptohopper response
//
//
//	response := `{"signal": "BUY"}`
//
//	assert.Equal(t, data, response)
//	assert.Nil(t, err)
//
//}

//func TestSignalSell() {
//	current_sma8 := 50
//	previous_sma8 := 49
//	sma55 := 40
//
//	response := `{"signal": "SELL"}`
//}
//
//func TestSignalNeutral() {
//
//	response := `{"signal": "NEUTRAL"}`
//}
