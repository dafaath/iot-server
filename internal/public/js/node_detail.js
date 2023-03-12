console.log("channel", CHANNEL);
for (let index = 0; index < CHANNEL.length; index++) {
  const options = {
    series: [CHANNEL[index]],
    chart: {
      id: `line-datetime-${index}`,
      type: "area",
      height: 350,
      zoom: {
        autoScaleYaxis: true,
      },
    },
    stroke: {
      curve: "smooth",
    },
    annotations: {
      yaxis: [
        {
          y: 30,
          borderColor: "#999",
          label: {
            show: true,
            text: "Support",
            style: {
              color: "#fff",
              background: "#00E396",
            },
          },
        },
      ],
      xaxis: [
        {
          // x: new Date("14 Nov 2012").getTime(),
          borderColor: "#999",
          yAxisIndex: 0,
          label: {
            show: true,
            text: "Rally",
            style: {
              color: "#fff",
              background: "#775DD0",
            },
          },
        },
      ],
    },
    dataLabels: {
      enabled: false,
    },
    markers: {
      size: 0,
      style: "hollow",
    },
    xaxis: {
      type: "datetime",
      // min: new Date("01 Mar 2012").getTime(),
      tickAmount: 6,
    },
    tooltip: {
      x: {
        format: "dd MMM yyyy",
      },
    },
    fill: {
      type: "gradient",
      gradient: {
        shadeIntensity: 1,
        opacityFrom: 0.7,
        opacityTo: 0.9,
        stops: [0, 100],
      },
    },
  };
  console.log(`#channel-chart-${index}`);
  const chart = new ApexCharts(
    document.querySelector(`#channel-chart-${index}`),
    options
  );
  chart.render();
}
