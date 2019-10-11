let counter = 0;
let data = [[]];
let funcId = 0;
let globalHeatLayer;
async function drawLayer(map, actuFuncID, operatorList) {
  var heatMapData = [];
  let arrayCoord = []
  operatorList.forEach((operator, index) => {
    let arrayUsed = []
    if (data[index][counter]) {
      arrayUsed = data[index][counter].scooter_list;
    }
    const cords = arrayUsed.map(scooter => [
      scooter.latitude,
      scooter.longitude,
      1
    ]);
    arrayCoord = arrayCoord.concat(cords)
  });
  let heatLayer = L.heatLayer(arrayCoord, { maxZoom: 18 });
  map.addLayer(heatLayer);
  globalHeatLayer = heatLayer;
  await sleep(1000);
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
  drawLayer(map, actuFuncID, operatorList);
}

async function loadData(operator) {
  const response = await fetch("http://localhost:3000/data?" + "operator=" + operator);
  const dataJson = await response.json();
  if (!dataJson) return []
  return dataJson;
}

async function refresh() {
  const operatorList = [];
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
  data = await Promise.all(operatorList.map(operator => loadData(operator)));
  counter = 0;
  drawLayer(globalMap, funcId, operatorList);
}
