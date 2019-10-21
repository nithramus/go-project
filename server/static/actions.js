let counter = 0;
let data = [[]];
let funcId = 0;
let globalHeatLayer;
let operatorList = [];
let map
let isInPause = false

function resume() {
   map.removeLayer(globalHeatLayer);
  funcId += 1
  isInPause = false
  counter = parseInt(timelapse.value);

  drawLayer(funcId)
}

function pause() {
  isInPause = true
}

const timelapse = document.getElementById("timelapse");
async function drawLayer(actuFuncID) {
  const date = new Date(data[0][counter].date)
  document.getElementById("infos").innerHTML = "Heure: " + date.getHours() + "h" + + date.getMinutes() 
  var heatMapData = [];
  let arrayCoord = [];
  operatorList.forEach((operator, index) => {
    let arrayUsed = [];
    if (data[index][counter]) {
      arrayUsed = data[index][counter].scooter_list;
    }
    const cords = arrayUsed.map(scooter => [
      scooter.latitude,
      scooter.longitude,
      1
    ]);
    arrayCoord = arrayCoord.concat(cords);
  });
  let heatLayer = L.heatLayer(arrayCoord, { maxZoom: 18 });
  map.addLayer(heatLayer);
  globalHeatLayer = heatLayer;
  await sleep(200);
  if (isInPause === true) {
    return 0;
  }
  timelapse.value = counter;
  if (counter + 1 >= data[0].length) {
    globalHeatLayer = heatLayer;
    return 0;
  }
  if (actuFuncID != funcId) {
    return 0;
  } else {
    map.removeLayer(heatLayer);
  }
  counter += 1;

  drawLayer(actuFuncID);
}

async function loadData(operator, date) {
  const response = await fetch(
    "/data?" + "operator=" + operator + "&date=" + date.toUTCString()
  );
  const dataJson = await response.json();
  if (!dataJson) return [];
  return dataJson;
}

async function refresh() {
  operatorList = []
  document.getElementById("infos").innerHTML = ""
  const lime = document.getElementById("lime");
  if (lime.checked) {
    operatorList.push("lime");
  }
  const bird = document.getElementById("bird");
  if (bird.checked) {
    operatorList.push("bird");
  }
  const hive = document.getElementById("hive");
  if (hive.checked) {
    operatorList.push("hive");
  }
  const circ = document.getElementById("circ");
  if (circ.checked) {
    operatorList.push("circ");
  }
  const tier = document.getElementById("tier");
  if (tier.checked) {
    operatorList.push("tier");
  }
  const voids = document.getElementById("void");
  if (voids.checked) {
    operatorList.push("voi");
  }
  const wind = document.getElementById("wind");
  if (wind.checked) {
    operatorList.push("wind");
  }
  if (globalHeatLayer) {
    globalMap.removeLayer(globalHeatLayer);
    globalHeatLayer = null;
  }
  funcId += 1;
  const date = new Date(document.getElementById("date").value)
  data = await Promise.all(operatorList.map(operator => loadData(operator, date)));
  if (data[0].length < 1) {
    document.getElementById("infos").innerHTML = "Aucune donnée disponible, essayez de change la date"
    return
  }
  // updating the timelapse
  timelapse.max = data[0].length - 1;
  timelapse.value = 0;
  counter = 0;
  map = globalMap
  drawLayer(funcId);
}
