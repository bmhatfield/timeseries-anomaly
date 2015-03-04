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
    var trainingpoints = makeWindowedSeries(series.trainingnpaa, series.paaWindow);
    var testingpoints = makeWindowedSeries(series.anomalynpaa, series.paaWindow);

    $.plot($("#discrete"),
        [
            {data: trainingpoints, lines: {steps: true}},
            {data: testingpoints, lines: {steps: true}}
        ],
        { series: {lines: {show: true}} }
    );

    $.plot($("#surprise"),
        [ {data: makeWindowedSeries(series.surprise, series.paaWindow * series.wordLength), lines: {steps: true}} ],
        { xaxis: {min: 0, max: trainingpoints.length * series.paaWindow}}
    );
});
