
function getRandomArbitrary(min, max) {
  return Math.random() * (max - min) + min;
}
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
function initMap() {
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
  document.getElementById('date').valueAsDate = new Date();
  
  
    setTimeout(function() {
        sidebar.show();
    }, 500);
  map.addLayer(Esri_WorldImagery);
  return map

}


// On initialise la latitude et la longitude de Paris (centre de la carte)
var lat = 48.852969;
var lon = 2.349903;
var macarte = null;
// Fonction d'initialisation de la carte

var globalMap

window.onload = async function() {
  // Fonction d'initialisation qui s'exécute lorsque le DOM est chargé
  globalMap = initMap();
  await sleep(2000)
  refresh()

};
