function getRandomArbitrary(min, max) {
  return Math.random() * (max - min) + min;
}
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
function initMap(data) {
  var Esri_WorldImagery = L.tileLayer(
    "http://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
    {
      attribution:
        "Tiles © Esri — Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community"
    }
  );

  var map = L.map("map").setView([48.852969, 2.349903], 13);
  var sidebar = L.control.sidebar("sidebar", {
    position: "left"
  });

  map.addControl(sidebar);
    setTimeout(function() {
        sidebar.show();
    }, 500);
  map.addLayer(Esri_WorldImagery);

  drawLayer(map, data, 0);
}
async function loadData() {
  const response = await fetch("http://localhost:3000/data");
  const dataJson = await response.json();
  return dataJson;
}
async function drawLayer(map, data, counter) {
  var heatMapData = [];
  const arrayUsed = data[counter].scooter_list;
  const arrayCoord = arrayUsed.map(scooter => [
    scooter.latitude,
    scooter.longitude,
    1
  ]);
  var heatLayer = L.heatLayer(arrayCoord, { maxZoom: 18 });
  map.addLayer(heatLayer);
  await sleep(1000);
  if (counter + 1 >= data.length) return 0;

  map.removeLayer(heatLayer);
  drawLayer(map, data, counter + 1);
}
// On initialise la latitude et la longitude de Paris (centre de la carte)
var lat = 48.852969;
var lon = 2.349903;
var macarte = null;
// Fonction d'initialisation de la carte

window.onload = async function() {
  // Fonction d'initialisation qui s'exécute lorsque le DOM est chargé
  const data = await loadData();
  initMap(data);
};
