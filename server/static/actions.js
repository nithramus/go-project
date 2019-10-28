let counter = 0;
let data = [[]];
let funcId = 0;
let globalHeatLayer;
let operatorList = [];
let map
let isInPause = false
let coordinateMap = []
let operatorMap = {}
function resume() {
  console.log(counter, data[0].listDiffs.length)
  if (counter + 1 >= data[0].listDiffs.length) {
    counter = 0
  }
  else {
    counter = parseInt(timelapse.value);
  }
   map.removeLayer(globalHeatLayer);
  funcId += 1
  isInPause = false

  drawLayer(funcId)
}

function pause() {
  isInPause = true
}

function createMap() {
  operatorList.forEach((operator, index) => {
    data[index].base.scooter_list.forEach(scooter => {
      operatorMap[scooter.id] = scooter
    })
  })
}

function updateOperatorMap() {
  operatorList.forEach((operator, index) => {
    if (data[index].listDiffs[counter].added) {

      data[index].listDiffs[counter].added.forEach(scooter => {
        operatorMap[scooter.id] = scooter
      })
    }
    if (data[index].listDiffs[counter].removed) {
      data[index].listDiffs[counter].removed.forEach(scooter => {
        delete operatorMap[scooter.id]
      })
    } 
  })
}

function createArrayCord() {
  if (counter === 0) {
    createMap()
  }
  else {
    updateOperatorMap()
  }
}

function createFinalArray() {
  return Object.entries(operatorMap).map(([clé, scooter]) => {
    return [
      scooter.lt,
      scooter.lg,
      1
    ]
  })
}

const timelapse = document.getElementById("timelapse");
async function drawLayer(actuFuncID) {
  const date = new Date(data[0].listDiffs[counter].date)
  document.getElementById("infos").innerHTML = "Heure: " + date.getHours() + "h" + + date.getMinutes() 
  createArrayCord();
  let arrayCoord = createFinalArray() 
  console.log(arrayCoord.length)
  let heatLayer = L.heatLayer(arrayCoord, { maxZoom: 18 });
  map.addLayer(heatLayer);
  timelapse.value = counter;
  globalHeatLayer = heatLayer;
  await sleep(200);
  if (isInPause === true) {
    return 0;
  }
  if (counter + 1 >= data[0].listDiffs.length) {
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
  document.getElementById("infos").innerHTML = "Chargement des données"
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
  // document.getElementById("infos").innerHTML = ""
  if (data[0].length < 1) {
    document.getElementById("infos").innerHTML = "Aucune donnée disponible, essayez de change la date"
    return
  }
  // updating the timelapse
  timelapse.max = data[0].listDiffs.length - 1;
  timelapse.value = 0;
  counter = 0;
  map = globalMap
  drawLayer(funcId);
}
