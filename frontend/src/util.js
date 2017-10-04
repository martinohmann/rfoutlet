import { api } from './config';

export function makeApiRequest(requestUri, data, success) {
  let formData = new FormData();

  for (var key in data) {
    if (data.hasOwnProperty(key)) {
      formData.append(key, data[key]);
    }
  }

  fetch(api.baseUri + requestUri, {
    method: "POST",
    body: formData,
  }).then(response => {
    return response.json();
  }).then(result => {
      return success(result);
  }).catch(err => {
    console.log(err);
  });
}

export function isOutletEnabled(outlet) {
  if (undefined === outlet || undefined === outlet.state) {
    return false;
  }

  return 1 === outlet.state;
}
