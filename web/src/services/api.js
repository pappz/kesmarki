import axios from 'axios';
import qs from 'qs';

const api = axios.create({
	baseURL: 'http://127.0.0.1:8180/',
	withCredentials: true,
	paramsSerializer: function(params) {
		return qs.stringify(params, { arrayFormat: 'repeat' });
	}
});

api.interceptors.response.use(
	function(response) {
		return response;
	},
	function(error) {
		if (error.response.status === 403 || error.response.status === 401) {
			window.location.href = `/logout`;
		}

		// Do something with response error
		return Promise.reject(error);
	}
);

export default {
	shutterUp() {
		return api.post('api/shutter/up');
	},
	shutterStop() {
		return api.post('api/shutter/stop');
	},
	shutterDown() {
		return api.post('api/shutter/down');
	},
}
