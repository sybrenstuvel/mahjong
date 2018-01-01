$(function() {
    toastr.options.closeButton = true;
    toastr.options.progressBar = true;
    toastr.options.positionClass = 'toast-bottom-left';
    toastr.options.hideMethod = 'slideUp';
})

function random_hand() {
    $.get('/api/random')
    .done(function(data) {
        console.log('hand randomised', data);
        $('#json_input').val(JSON.stringify(data, undefined, 4));
    })
    .fail(function(err) {
        toastr.error(err.statusText, 'Unable to get random hand');
    })
    ;
}

function score_hand() {
    $.post('/api/calc-score', json=$('#json_input').val())
    .done(function(data) {
        console.log('hand scored', data);
        toastr.success(data.score, 'Calculated score');
    })
    .fail(function(err) {
        toastr.error(err.statusText, 'Unable to score hand');
    })
    ;
}
