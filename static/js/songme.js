
$('.embed-link').each(function(){
    link = $(this);
    url = link.attr('href');
    $.getJSON('https://noembed.com/embed', {format: 'json', url: url}, function (response) {
      link.replaceWith(response.html);
    }).fail(function (){
        console.log('Can not fetch embed code');
    });
});