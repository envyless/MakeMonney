package scrape_define

//scrape은 아무래도 웹사이트에 의존하게 된다. 웹사이트의 특성상 구조가 변경될 수 있다.
//그렇게 된다면 변수취급 되어야 하며, 그런 것들은 조금은 정리해서 설명과 함께 한곳에 모아놓는 것이 경험상 편하더라.

//모든 함수키를 여기서 만든다.

var MainData = "section inner_sub"
var FindKey_Daily_Price = "section inner_sub"

//일별시세 관련 데이터

//날짜, 종가, 전일비, 시가, 고가, 저가, 거래량
var FK_date = "날짜"
var FK_lastPrice = "종가"
var FK_deferPrice = "전일비"
var FK_cap = "시가"
var FK_highCap = "고가"
var FK_lowCap = "저가"
var FK_volume = "거래량"
