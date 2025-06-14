package storage_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"done-hub/common/utils"

	"done-hub/common/requester"
	"done-hub/common/storage/drives"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testImageB64 = `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAABg8SURBVHgB5V0JmBXVlf6r6vXrjV5EEBC1H3FBCQpk4qdxxZiMfuNkkIlLtk8xe5wAmhj3ETAhRJO4TFRc4odOMJtxhKjopyItKkExsqgIsnSziPS+vn5LLXfOuVWvebRvqXpVhZj8evu9bmq5dc89+7m3FBzkEELU08dUajFqDc4nt/qslo1upzU7v6+jtsP5XKcoSjcOYig4yEAEiNHHhdQmYR8hgkQzbOIshU2gdTiIcFAQhIgwlT6mwSZEDAcWzdQaqT1KxGnEPytYFFGbQ22tOHjQRG2GsLn0nwP0sJOp3UmtSwQAy/mR7uwQpmGKALFI/CMThh+O2goRMCz+zzRF+ykTRXrTRhECFol/JMIIWzTdKcICcQcTpHVcnWifdKRIbX5PGIYuQsAi8UknDD3AbBGQaMoLhyAdsXrRMkYVbRMPF+kNG0RIaKI2A580iJDEU04Mcki9aBujSKK0Txgjkm9vEIZpiJCwQoTELSoCBnV0Nn2she1DHHCoQsDsbkH8onNhrluPkDCVGluHVyFgBEYQ4egK+noXPuo9HzgogohiwerrQO/Xz0fqrTdhmRZCAD/jnc4zB4ZACOKw7wpqgc+YUiEs8nl7O9B36XlIr1mNEHGVsHVLDAHAN0GyiDEZIcCyTJjUvEJRTKYKRLwH/ZddgNTqV+W1QkKMWiB6xRdBqANMhLUIMdzRlABW9fgYSNIpor8P/V//N6RfbUSIiMHWK74mZskEcW7MnBGKvtCFiffjOmZvFoj6iLgpikUBOxMiEUf/jOlIvvwScZyBkMBjscIPUUoiSNjEYGyPm7hmE9BjmDiy3L+qI+EFJIkol0+D/tLzCBG+iOL5SbN0RjicYel4t0/HVZsEunULEepifVSDX6hMECKLpQ8g/u2LkXx+GciFQUjIECUGj/BEkLCJwdjcr+BaElPEICRuNERI5ETV4AZOISJY6ST6v3sx0sv+ihBRElG8cggTI4YQQCFCrOsx8KMtOuIGJ2pYcVjQ6MM/f+yDQpyisABLJ9B35aVILFtC1lcofgojRu1JYWc9XcE1QRwHKIaQ8FaPheu2WEgaGqysvJkp7BY4FCK5nkb/976K1JK/IESwLpnj9mBXBHHCIaE4fWT/4PVuA9dvMTBg2ClMRdlHgTR9TRWYwIJ1g7BDJkK4N8cUYcmzYaXQN/syJJ74Y5jW11VuwyxFCeLIwLkICW90mbhpqwnd1GTY4yMgj7vLyO+HKDzTqYnyitIS0nRLxUojPvtypP+wGCFijht94oZDFiEEJc5y+7VOEzcTMdKmYsegcowok+jDZJFuqgoio8Y4R3uD1CnMXWTdxa/7LhKLF0l9FgJ4DBcVO6jgkwo79j8VIWBVt8DcbcQZHHMqMJD8zzsTxQdaHX8CXaZ0D1KaEESU/uu+j+TDDyIkTC0muvISxGEv18rILXj2NXYauGWbjqRl21JKAVnDXPNevPiMjZx0slTUpUPIwVAE6bJb/guJ394XlvU1p5DVVYhDmBgxBIwVHSZ+utWAITQpLoqCxMkGFwTRzjwXQaLvltlILAw0sp4BEyPvRM85pRzuaEKAoLweXugQuI3FlPQFvM3mx05SMYZCKJqaew5xRLjn1BNg7d6CYGDfp/qG21D+w6tJTQXpDUmMI2OkOfddP4rARdWzRIxfkJjSJSG8x6Ze697fHB4KtrSiF34FQaN3wfVI/HoBQkBOBf+RaRo0d5gkh59pt3BHswU/oaPxwxQ8cAKJuTwcYvF/O3ai68wJdNM06YJgcx+VV9+Kyh/fCFUNlFOmDC1lzfV0gXLHU22mb2Iw3ic90lzA2uIMoXpkDJXTLkUYiN8xF/EFgQuOGUP/sB+HBMkd7PX+X4vAvbsGrw1foJ7+52HADxsittBTcnNKumkb+s+ZTAHEfoSByh/chIob50LTIggAXIk/Lrsif+hTTUVAWNoqcM8uUxJCBBDm5ks8R3ooTiH5QsZZpCGGsiuvRVgYWLgAiQVzERDY4trPLxlKEN88aZLsXtlh4Dc7hbSlggJfacBQ8Hgr/2IVOE5FxeyfQDvxVMezCBYcKU7eNx8DDwXmp0zL/mWwx8JeEhCDTzT1C/xsO5m25DUrIvgBeWKvhR49/3VNGrJIWTlqHngM6vDDEBbi834MvXE5AsBkZ+wlsp/scviASX5Gn67jv7cJpCg2pUm5EnzcvJ8Csr/bY0qPP1cVSYSeSKHYlnJUDDUPPk5CvzYcTqGAZP/My5DetZP64TtKfGHmS3ZPp8IH2Ou+c4eFPSnKyAX//Nk3wpJWk3LuRsHoLndBO/VzqFm0BEp1LYIGx9isrlYMzLoCwn8h3qDYkkPnJORjKAGybopm6yudFl7scC4aWqrapoFJovCOHQpZcvkpwv6KRj5D5MyzMOxPy6GMaaApE5wPwdWRrE/01S8j8fBCmUTzkaOPZULzmblcei0Rechs+dy1I0Qq5MC7pKueaCnu/LEei06ehLplq1F29gUIGoKiB8lfzYHY0STDQz4gxVaGINNQItik/R0p2k4dnuNTfsCR9od3W9hMSqVQ/kIlfcLedWTUaNT9fgmq7vwdMHystACD4BiVuAQDPeiff7Pf559kX89GyRyyl3KsT7bggENYZDzQ7J9Hwcq44W5mcryr4qKvov7l9aj87rVQKochKBjL/gzjzTfgA1P5h+LE5rvgEZYlZD77V80GlrXhYwPn0c8YbuLWo8tIb9CcV9x50Cb3v70Nqf99CInfP0RpyV0Qqp0QUNm/ULyKYOLEM89H3R+flm6SqpVk2RyiODbwCniEIAewhSyqr71tgSfox7a+mp5eWBq+cTjwrSMipMjd9cRyBl0hglomRaFXNiK99M9Iv/g0rO42p7DOA4hbBXFg7VOroU2ZQuZ3SeJwOk+nksQVC4mlbRZMn8SQVSNcqFCqTUADwWH5xz4ERkcNXHCYJkWTWiR7qGaixlz3pZZDO/c8VFAz0ykSPW9Cf20FzDWvIb12DZSBTqeiRZGWlcxMDtVbil28lPrtPai+92GUiBhzCC+wmQ2PSJsmLllvokuHL7ByrdMM9Jo+FSzTVbMw92gFZ9WX7Rtwr5cR2ZFpGmIjDXPD29BfWY7US8/CWLcKimkgr9MbrUbdmu0oG1lSlOBu7vUklIA1PUCnDwdVCJsYJ9cBtx9fJkWEL8NZTloV87eqWNVjSJFklmCGKhRFZmLajTgtWo6yz34WVVdfh0OeasQhrzej8qZfQo1NkOJOKNp+MTsrnYCxtOTCu1jJPvXyTsuXqGKur44IXPspBcdWqfjSYcwh/jVRiogwZ4uC1UQURQSv2dTRo1Bx5dWoW7ke1fc/ToQ5fr9/Z6GW+OvjKBENTJCYlzMsYUAnxbG6W/hyhFgmf2ssMCoaIQWo4ttHAiPL/Gf5eEAMYr+bKbW+vMuQkQQvsSZBx3J+no0WbpnvGbBPoxEXaZEIKqZ9GfUvvoHK2XPpHyoonKJKHWOuXQWjtUVach5RzwTxVATHnu+7lL3r4+I2HzPw8ApBXLGPQauJKD+KBRfaIP7Az7YI/KXF9CQKeYr9+P5WXPHLFtz/dDe27TUKpxHKK1D1k5tQ89gzUGuGyz8pZLUZK18qWAOQB94JwrdY12evt1BKqIPK1N9eOlpBWVYlB1eTnD68DBeOYi6Eb3A1Io/jvTuB3+zigjyOEBfnQLbOzju5Gs+93od5i9px1qwmfOHaHXhoWQd6+nV5jewosyZ1DcfMzqaY2QtA7Qhp/uqrVqIE1HvXIXSzjf2ljxjTsIKsoS+OyF2xeCX5EscPC1L2U8xrL3HLVqtgMHLf0Qq+8JkqjBlpp4rZrH5vm46bf9uOz/6gGT//Qye647kIqyB64okY9uCfyEGNQn9rdUlGimeCsC+7ecDfFD69XsMwksG58uIRsmxuPVaAmCWgbIoiH7KRYhF/bLHjXgVjX9SnMsqXTz+zzr4/m8FygYRAHxHi3r904PSZO/DoC50wTNPRUWKQU8rOOAsVs26CuW0zZbES8ArPBOkhv4OXmpUKfsgz6guUjpLcHRUtw0+PVVEZaBxf4JHdwK6k5api8jziEpFDJHMStLvXxHUL2/GV+XvQ0mUMSSkrqJx1DbTY0TC3b4dXMEE87UG4N8WKr/QMFK+ImliX/3zpB9BATKhWMecYSsf6sx32g0F0eGyPgBuDfcqxUVRU5OinU2PBCvu19QM474bdWLct5TiUackpGiv6626F2O65irLbO0Eouqv4CPuPKBPUih/HRDi1LoIbyU+JBMgoKzrs9YvFEKHg4KcbokWPa2s38eW5H+C19+LU6X2BzfLzL4B2zhfgEd4JItfw+5ixR1S6m/Ka4zFPPVTFjcdwbAp+fXmJJN17Q0/xGcU9PPbw/ATJlDcxvyUGTHxzwV68szvl6BTiwUgUkTrPqeNmJkizlzPk7PJRZ3WYnETuz+cOnjOc8h5ElHLFf7Ke1cLGAVdHYtQI9/frJUPnO7d9iN6E7ivewHfc4eUE3fJeuZ6NYREl50qpfOCl0Rq1M+oj+MV4O9zip96La37bdVcHorYqktUPpaDfxYbCjr06fv5Yl0zrloj1njnEN5TBHx5hYUpNBPedoGFMuR8rj4Iboa2CBha/0IfNH6RQIryLLJm78MGThmXLXa9gG5+jr0dVqlg4QcWUOqU0R4UdUxeSiPWDmVXeM7QkNh/HWKbA/U91O1aX59jcOukzeTmjUvGnWntNxZfI47vXRzTcTn7KWYfCO0ipH1FZ/DAWrB09JTwpjc9Tq/qRkL6aZ523TnUqr5vdnlETUeCH4/ckLXsjmBLB+kSVNVcKftRAil61cjpweUEEmegiNMMzfHdHepAThnLEUI4Z/J3+H6Dg65rNSTue5h7NTIsMCV92e1Z9xPI1w3emfBlpg+CO15ZFMLlW9VRdP7pc4Pjq4v3nwd+0o/R06HpyFi1vwyQ3iMwQZJ3bs0aXq65CD/kQJ3e5OeGfIuzRazRoVbK6o/iTZ4Lw00e5639Hv4nmD1MlL6fY2eLZ/F3CP9TsX9yACVLmZ0MxEhlv9gXjevNAbY6brgaY1/yOKQcuPExxlTZ49d2Eq+hwPsRTnm1JyRSSIM5q0GY3Z5UrJj2YD4pQwHAlhS9kxLWEmWd7x4KsNQvPtumkkwoL0EwNvkr3vfFTZJSQQaAWcDAtJ1O47G99+1XCZDilmD+SgdxxyL1ua86sNczu2VI3Z/I9jvdT8Eez7t04OVFJq7QtLBTbOX2mVcevXbi0PKYcOrt5nIYJw1wU0REHd/YZeG5NvKQEXAbDh2lezPLGzJfsHrLYclEOpGJijcDz7SXmvxWbMZa00M0a3Pe4LZnGr3cqSBMxtg8o6NILD5bt66hkhAA3UIDytEPcVTQy9y1e3o90MvdkyeiTDLHy6ZcGGQeTcWG4wKOZL4Mc4rzMpKhy5/v/S61/V/eZVoG2lHs7fXi5hhtiChrIq+s3ik8GEkz41xHAwxM1nFLnfoFmX1Lgwac64RcTG8rdpqKbs18kM7SnLLYKVjKyD3BEBTcdu5OlK+cUnfrAbgM3HR2RefZIkRJQ9j/qaNLNbDDx9TEaVtCYrekli42i3p1OwLOOLK6jqoHP1ACfJ6dxTHlhfZENjtByF+56ohPtPVbRysd8nCHLm6o1TDkmCpdVrfP2O3/ITbjgoQkuCh8e2Glg8Yemp0Dhfvdybr7gOIXyHprrSkMuflNlgZpdwyoz85YqnU1pADvdkQJLIO9GA0PBBNnQPIB/v34PmAGLWW6FTOHpU2tx38zRcqKpxamy3xYb+/XW8dofhQuce6giTclSkTnz9iaBDl04awZdmK9c+8QVheyxyx7YMa6IU23I/onto6iuiJGxnvqSJq68sw0GV47LOJQo2PJej2bE9DNqnArIouPzyND9TnL1+BG4wDiKBx0XQHVINznDt2y1kDJLWQIQDExyVmfevRdNe9LwizEjIph6UpXbw+cN/cNHCOLYw40oArYypjv1xH6Gkc/lsiJevZt21ulZ4e19OAgWURyv0k0TVz2wF8+TmZtr5hfNgwzGuVQpL7/xxXqUuZOSj3jZDegKuMDnhysYERWBpFbf6CbH7X0TcZ0FuP/MYDHwGCfpXjPvbcUTK3rhG6TTqsiImPHFWrdFGfNy/THnkzuUuxsFwNZLuRbBV0YrgRQ18xX+TuHuKzdZ2N4nPNfkugVfk/XVBx06Lrr1Qyx9ubekYOfg5ptZUd5vnF+LEXVlbvbWyskd8rr5znBrcSVo4Ga8bWJvEoGAraVKerjL60xcfEw5omXBbhymmwaefLUftzzSju5eo6i8zecADv17DZm6q+6JYWRd0ZKaZmrn5COIWqAjbHHNQxFUUL++M1aRJl4Q9VNsRqfoYRf2abjs2X68uC5OcatMxWEJ0QHLzr/w+X/b1I+L536AWf/Tiu4ew5Xyy2dVDVpcDmFmXXQIDq1xNQDz8hGDUfQK/AIsFNjlQb5shTr1k/fT+HuXCgRAFBscnFOht6RQt7EHl0yqxJdOq8XIWhdFXVnoiafxLCnsxS/24c1NCXsLBvmyl2A6yjP66KMq8cLtYxGNqMX8qUYixjmFDnBDkBjsl7YUFF0fJA188x0LyRBeYsMzsevNXvSs6sGkT5XjlAlRnEh297ix5Rhdp6IiGpGVhKm0ibZeymPs0fHOjiTe2JjC68QV6bQSTFYsB1RNw5PzD8cpx1UXO5QlzpRC3MFwNU2cvWYLbtHJ4fDnOgz80ns5qyuw4WAMGOh4pRPdbw84itTmIk21Y2u85Ybcchy2l57ZEthWvuGUmsy+5FBcf8khZBgWjZddTf24q9hBrvnW7eLQW7emsbxdCc3JY0mjd6bRuaobfe8RYTLSp8iAF4vOuoV9HTvwc9rJ1fjDNYcjygXLhWNfd9N5rvZ+90IQ+T4MFAk+DlDoYdZ7Bra4qg70DktWVqmS4Fafhe613ehdT2IpXthEDo4gdvr4+DNr8fj3RuLQcnu5qpLfd+Jo7ji4hCfNlvVCl1i+Yzge1arrmLVRQUva8nqLksCDnGhKSI7p29ZHok2RGcLsdITI/OKWIM55to/lnE2DXju+Cg1n12HhGVGMq4oWiwo3o4CJW+C27lFMyWe2KNqVIE7ZZMpY1YGAHe21ze/k3gSSu5JIkHJPk5Wm9+qygE2u93d5vczAaDURVI6tRPW4Cgw7jj6rIrhtvIaTauyLFbCqXCnxfPf1BDcvBWO7f3u/jmveV9GpCyjhM0pOCOYUsrL03jT0HvL+4xTIpH5pFFAUukUWmK17IhUaNN6OLkqRY0q/RijVGK0vR6RKdfL4CqrpWvPHC0weFi0WyWVinDN0T143KHmYihHF3kuEzOGEiWu2iMA8ea+wZE5EyKVP1uBSI3srD5lByVrsIrWTs1jULo7Y9woNTgXPH6/i01Wa5EatMGeURAy7Zz7gllM60gbmbgXeCWcr3VDBiTBOG//8OGBseaRYBtIXMRi+BUkxRS+cbZx0mnn37DTx11Z8onD2cODacaosyGMFr+QXVc3wqMBzIRDJ7sb6ymBlp447miy5aQ3vT/Ux5aQKgrtUrQn84EjgP0aVuzmFOWK6X2IwAlW1bpxH3py/hxTqfcQtL3YEu9lyEODK5dPrLcw8KgKmhYt9rzhNMTd7u3A/CHw0nDDLHBSJfXF14NuU91i428K7vRnP3t/qrFKQbS43VAt8/wgVn6tT3LwFQUbD3YRDvCCUp3dE2CIUiBLLBS2wxcOaLoHFewxsiPsr5C4FbFQdQ3HBr41RSF9octm2tL8KJ5kaqV0RhIgailCno/NSMeaWWLFj2UzeSD7C022kZzos9PNKVt7zk4wa1fQfG5MhE8URkWQpVdC304ar+NJIgSk19h6/LhAKV2QjdPmQ9R7Ey4sc5ywBU8GrFdZ0W3iF8uxrqbUb/jtq8X6IpKin1Ko4i4TpafUqKrXMksfiWwIiYF2RDwdMYLslzFBwYdyuBPswggKWJn1XsCtNiaeUQFI6eJbcp4rBUXjmKJ79dWWKXJzTwOVK5G2fMExFjL5HvW/914iQxFMuHHATxythBjlHzmB7gzDutklKOG5alNPn77Y447K5CvLKq8tURGBltqx0wor2T5elpZmCwUf8OHmfKDBhWMdQaxIHD9ZSmys8vN05aBwUToCw9w6eQe1suDAAAkYz7CLzJdlV6B8XDi6vDIPxMW7TnM8YgkUzbL3AiyyXHCjd4BYHHUGGwhEfGSLFsG9bW/5ej486oN1ZrRn21iHNTmsM20ryi/8HGjuXYi/dS74AAAAASUVORK5CYII=`

func InitConfig() {
	viper.AddConfigPath("/done-hub")
	viper.SetConfigName("config")
	viper.ReadInConfig()
	requester.InitHttpClient()
}
func TestALIOSSUpload(t *testing.T) {
	InitConfig()
	endpoint := viper.GetString("storage.alioss.endpoint")
	accessKeyId := viper.GetString("storage.alioss.accessKeyId")
	accessKeySecret := viper.GetString("storage.alioss.accessKeySecret")
	bucketName := viper.GetString("storage.alioss.bucketName")
	aliUpload := drives.NewAliOSSUpload(endpoint, accessKeyId, accessKeySecret, bucketName)

	image, err := base64.StdEncoding.DecodeString(testImageB64)
	if err != nil {
		fmt.Println(err)
	}

	url, err := aliUpload.Upload(image, utils.GetUUID()+".png")
	fmt.Println(url)
	fmt.Println(err)
	assert.Nil(t, err)
}
func TestSMMSUpload(t *testing.T) {
	InitConfig()
	smSecret := viper.GetString("storage.smms.secret")
	smUpload := drives.NewSMUpload(smSecret)

	image, err := base64.StdEncoding.DecodeString(testImageB64)
	if err != nil {
		fmt.Println(err)
	}

	url, err := smUpload.Upload(image, utils.GetUUID()+".png")
	fmt.Println(url)
	fmt.Println(err)
	assert.Nil(t, err)
}

func TestImgurUpload(t *testing.T) {
	InitConfig()
	imgurClientId := viper.GetString("storage.imgur.client_id")
	imgurUpload := drives.NewImgurUpload(imgurClientId)

	image, err := base64.StdEncoding.DecodeString(testImageB64)
	if err != nil {
		fmt.Println(err)
	}

	url, err := imgurUpload.Upload(image, utils.GetUUID()+".png")
	fmt.Println(url)
	fmt.Println(err)
	assert.Nil(t, err)
}
