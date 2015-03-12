package engine

const (
	DASHBOARD_TPL = `
<html>
<head>
<title>{{ .Title }}</title>
<meta http-equiv="refresh" content="10">
<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.time.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.selection.min.js"></script>

<script type="text/javascript">
	var calls = {{ .Calls }};
	var sessions = {{ .Sessions }};
	var data = [
			{ label: "slow", data: {{ .Slows }} },	
    		{ label: "conns", data: {{ .ActiveSessions }} },
    		{ label: "err", data: {{ .Errors }} },
    		{ label: "qps", data: {{ .Qps }} },
    		{ label: "latency", data: {{ .Latencies }} },
	];

	var options = {
		legend: {
			position: "nw",
			noColumns: 5,
			backgroundOpacity: 0.2
		},
		xaxis: {
			mode: "time",
			timezone: "browser",
			timeformat: "%H:%M:%S "
		},
		selection: {
			mode: "x"
		},
	};

	$(document).ready(function() {

	var plot = $.plot("#placeholder", data, options);

	var overview = $.plot("#overview", data, {
		legend: { show: false},
		series: {
			lines: {
				show: true,
				lineWidth: 1
			},
			shadowSize: 0
		},
		xaxis: {
			ticks: [],
			mode: "time"
		},
		yaxis: {
			ticks: [],
			min: 0,
			autoscaleMargin: 0.1
		},
		selection: {
			mode: "x"
		}
	});

	// now connect the two
	$("#placeholder").bind("plotselected", function (event, ranges) {
		// do the zooming
		$.each(plot.getXAxes(), function(_, axis) {
			var opts = axis.options;
			opts.min = ranges.xaxis.from;
			opts.max = ranges.xaxis.to;
		});
		plot.setupGrid();
		plot.draw();
		plot.clearSelection();

		// don't fire event on the overview to prevent eternal loop

		overview.setSelection(ranges, true);
	});

	$("#overview").bind("plotselected", function (event, ranges) {
		plot.setSelection(ranges);
	});

	});
</script>
<style>
#content {
	margin: 0 auto;
	padding: 10px;
}

.demo-container {
	box-sizing: border-box;
	width: 1200px;
	height: 450px;
	padding: 20px 15px 15px 15px;
	margin: 15px auto 30px auto;
	border: 1px solid #ddd;
	background: #fff;
	background: linear-gradient(#f6f6f6 0, #fff 50px);
	background: -o-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -ms-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -moz-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -webkit-linear-gradient(#f6f6f6 0, #fff 50px);
	box-shadow: 0 3px 10px rgba(0,0,0,0.15);
	-o-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-ms-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-moz-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-webkit-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.demo-placeholder {
	width: 100%;
	height: 100%;
	font-size: 14px;
	line-height: 1.2em;
}
</style>
</head>
<body>
<pre>{{ .Title }}</pre>
<div id="peers">
  <p>Peers in cluster</p>
  {{range $index, $peer := .Peers}}<a href="http://{{$peer}}">{{$peer}}</a>&nbsp;&nbsp;
  {{end}}
</div>

<div id="content">
	<div class="demo-container">
		<div id="placeholder" class="demo-placeholder"></div>
	</div>

	<div class="demo-container" style="height:150px;">
		<div id="overview" class="demo-placeholder"></div>
	</div>
</div>

</body>
</html>
	`
)
