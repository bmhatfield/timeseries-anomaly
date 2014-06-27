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
    var trainingpoints = makeWindowedSeries(series.trainingnpaa, paaWindow);
    var testingpoints = makeWindowedSeries(series.anomalynpaa, paaWindow);

    $.plot($("#discrete"),
        [
            {data: trainingpoints, lines: {steps: true}},
            {data: testingpoints, lines: {steps: true}},
            {data: [[0, series.breakpoints[0]], [trainingLen-1, series.breakpoints[0]]], color: 'black', shadowSize: 0},
            {data: [[0, series.breakpoints[1]], [trainingLen-1, series.breakpoints[1]]], color: 'black', shadowSize: 0},
            {data: [[0, series.breakpoints[2]], [trainingLen-1, series.breakpoints[2]]], color: 'black', shadowSize: 0},
            {data: [[0, series.breakpoints[3]], [trainingLen-1, series.breakpoints[3]]], color: 'black', shadowSize: 0}
        ],
        { series: {lines: {show: true}} });

    $.plot($("#surprise"),
        [ {data: makeWindowedSeries(series.deltas, paaWindow), lines: {steps: true}} ],
        { xaxis: {min: 0, max: trainingpoints.length * paaWindow}});
});
