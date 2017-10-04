const apiBaseUri = 'http://127.0.0.1:3334/api';

function makeApiRequest(requestUri, data, success) {
  let formData = new FormData();

  for (var key in data) {
    if (data.hasOwnProperty(key)) {
      formData.append(key, data[key]);
    }
  }

  fetch(`${apiBaseUri}${requestUri}`, {
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

export default makeApiRequest;
