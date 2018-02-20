
// Runs when document ready.
$( document ).ready(function() {
    replaceEmbedCodes();
});

// replaceEmbedCodes finds each anchor element with class embed-link
// and requests for embed codes from url that anchor element points.
function replaceEmbedCodes() {
    $('a.embed-link').each(function(index, element){
        $.getJSON('https://noembed.com/embed', {format: 'json', url: element.href}, function (response) {
          $(element).replaceWith(response.html);
        }).fail(function (){
            console.log('Can not fetch embed code');
        });
    });
}
