
$('a.embed-link').each(function(index, element){
    $.getJSON('https://noembed.com/embed', {format: 'json', url: element.href}, function (response) {
      $(element).replaceWith(response.html);
    }).fail(function (){
        console.log('Can not fetch embed code');
    });
});