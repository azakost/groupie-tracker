// Run function at loading
var suggest = []
var pg = 1
updateList()

$('#search').on('keyup input paste keydown', function () {
    updateList()
})

$('#saveFilter').on('click', function () {
    updateList()
})

$('#clearFilters').on('click', function () {
    $('#modal input').val('');
    updateList()
})


function getPage(n) {
    pg = n
    updateList()
}

function updateList() {
    $.post("/api", {
        txt: $('#search').val(),
        page: pg,
        // sort: document.forms['filter'].elements['sort'].value,
        bycreation: $('#est_start').val() + '~' + $('#est_end').val(),
        bymembers: $('#mem_start').val() + '~' + $('#mem_end').val(),
        byalbum: $('#alb_start').val() + '~' + $('#alb_end').val(),
        byevent: $('#con_start').val() + '~' + $('#con_end').val(),
        bylocation: $('#loc_patern').val()
    },
        function (r) {

            // Making up paggination
            if (r.Found > 8) {
                $('#pag').empty()
                var l = r.Pages
                var p = r.Page
                var z = 1
                if (p != 1) {
                    z = p - 1
                }
                $('#pag').append('<li class="page-item"><a class="page-link" href="#" onclick=getPage(' + z + ')><span>&laquo;</span></a></li>')
                for (x = 1; x <= l; x++) {
                    cl = 'active'
                    if (x != p) {
                        cl = ''
                    }
                    $('#pag').append('<li class="page-item ' + cl + '"><a class="page-link" href="#" onclick=getPage(' + x + ')>' + x + '</a></li>')
                }
                k = l
                if (p != l) {
                    k = p + 1
                }
                $('#pag').append('<li class="page-item"><a class="page-link" href="#" onclick=getPage(' + k + ')><span>&raquo;</span></a></li>')
            } else {
                $('#pag').empty()
            }

            suggest = r.Suggestions
            $('#found').text(r.Found);
            // Appending cards template
            $('#cards').empty();
            $.each(r.List, function () {
                var tmpl = '<div class="col-sm-12 col-md-6 col-lg-4 col-xl-3 mb-4">' +
                    '<div class="card">' +
                    '<img class="card-img-top" src="' + this.Image + '">' +
                    '<div class="card-footer">' +
                    '<h5 class="card-title mt-3 mb-3">' + this.Name + '</h5>' +
                    '<div class="row">' +
                    '<div class="col-6">' +
                    '<p>Established</p>' +
                    '<p>Members</p>' +
                    '<p>First release</p>' +
                    '</div>' +
                    '<div class="col-6 rightflow">' +
                    '<p>' + this.Establish + '</p>' +
                    '<p>' + this.MembersNum + '</p>' +
                    '<p>' + this.Album + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</div>' +
                    '<div class="card-footer d-flex justify-content-center">' +
                    '<i class="material-icons mr-3 ' + this.Match.name + '">people</i>' +
                    '<i class="material-icons mr-3 ' + this.Match.member + '">person</i>' +
                    '<i class="material-icons mr-3 ' + this.Match.location + '">location_on</i>' +
                    '<i class="material-icons mr-3 ' + this.Match.album + '">library_music</i>' +
                    '<i class="material-icons mr-3 ' + this.Match.establish + '">cake</i>' +
                    '<i class="material-icons ' + this.Match.event + '">event</i>' +
                    '</div>' +
                    '<a href="/' + this.Id + '"></a>' +
                    '</div>' +
                    '</div>';
                $('#cards').append(tmpl)
            })
        }
    )
}









