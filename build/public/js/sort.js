$(document).on("click", '.sorterby', function(e) {
    if ($(this).hasClass('active')) {
        $(this).attr('data-sorter-direction', $(this).attr('data-sorter-direction') === 'asc' ? 'desc' : 'asc');
    }
    $(".sorterby").removeClass("active");
    $(this).addClass("active");
    $(".table-sort").change();
});