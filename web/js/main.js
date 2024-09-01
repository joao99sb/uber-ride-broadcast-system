let map;
let marker, circle;
let currentSocket = null;
const DEFAULT_LAT = -8.037544
const DEFAULT_LNG = -34.873450
const DEFAULT_ZOOM = 18;

function initMap() {
  map = L.map('map').setView([DEFAULT_LAT, DEFAULT_LNG], DEFAULT_ZOOM);
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map);
}

function updateMarkerPosition(lat, lng, zoom = DEFAULT_ZOOM) {
  if (marker) {
    marker.setLatLng([lat, lng]);
    circle.setLatLng([lat, lng]);
  } else {
    marker = L.marker([lat, lng]).addTo(map);
    circle = L.circle([lat, lng], { radius: 32 }).addTo(map);
  }

  map.setView([lat, lng], zoom);
}


function configSocket() {

  if (currentSocket) {
    currentSocket.close();
  }

  currentSocket = io("localhost:3001");

  currentSocket.on("travel_id", (data) => {
    const dataParsed = JSON.parse(data);
    const lat = dataParsed.lat;
    const lng = dataParsed.lng;
    if (lat && lng) {
      updateMarkerPosition(lat, lng);
    }

    if (lat == 0 && lng == 0) {
      alert("Viagem finalizada");
      currentSocket.disconnect();
      updateMarkerPosition(DEFAULT_LAT, DEFAULT_LNG);
    }

  });
}

function main() {
  initMap();


  const input = document.getElementById("input");

  document.addEventListener("submit", (e) => {
    e.preventDefault();

    configSocket();
    if (input.value) {
      currentSocket.emit("travel_id", input.value);
      input.value = "";
    }
  })


}

main();
