
var coords = [];
var art = window.location.pathname.split('/')[1]

$.post("/map",
    {
        id: art - 1
    },
    function (res) {
        coords = res
    }
);

ymaps.ready(init);
function init() {
    var myMap = new ymaps.Map("map", {
        center: [37, 20],
        zoom: 2,
        controls: []
    })
    $.each(coords, function () {
        var dates = ""
        $.each(this.Dates, function () {
            dates += this + "<br>"
        })
        myMap.geoObjects.add(new ymaps.Placemark(this.Coords, {
            balloonContent: "<b>" + this.Place + "</b><br>" + dates
        }))
    })
}

