function makeSeries(list){
    return $.map(list, function(i, index){
        return [[index, i]];
    });
}

function makeWindowedSeries(list, windw){
    return $.map(list, function(i, index){
        return [[index * windw -1, i]];
    });
}

$(document).ready(function(){
    $.plot($("#placeholder"),
        [
            makeSeries(series.normal),
            {data: makeWindowedSeries(series.npaa, paaWindow), lines: {steps: true}},
            [[0, series.breakpoints[0]], [rawLength-1, series.breakpoints[0]]],
            [[0, series.breakpoints[1]], [rawLength-1, series.breakpoints[1]]],
            [[0, series.breakpoints[2]], [rawLength-1, series.breakpoints[2]]],
            [[0, series.breakpoints[3]], [rawLength-1, series.breakpoints[3]]]
        ],
        {series: {lines: {show: true}}});
});
