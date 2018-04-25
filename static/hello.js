function OS($scope, $http) {
	$http.get('/os').
	success(function(data) {
		$scope.os = data;
	});
}
function Bookmarks($scope, $http) {
	$http.get('/bookmark').
	success(function(data) {
		$scope.bookmarks = data;
	});
}

var listBookmarks = angular.module("listBookmarks", []);
listBookmarks.controller("Bookmarks", ['$scope', '$http', function($scope, $http) {
	$http.get('/bookmark').
	success(function(data) {
		$scope.bookmarks = data;
	});
}]);
listBookmarks.controller("OS", ['$scope', '$http', function($scope, $http) {
	$http.get('/os').
	success(function(data) {
		$scope.os = data;
	});
}]);

var submitExample = angular.module("submitExample", [])
.controller('NewBookmarkController', ['$scope', '$http', function($scope, $http) {
	$scope.url = 'http://google.com';
	$scope.submit = function() {
		if(!$scope.url) {
			return
		}
		var body = JSON.stringify({URL:$scope.url})
		$http.post("/bookmark", body).
			success(function(data, status, headers, config) {
				// this callback will be called asynchronously
				// when the response is available
				console.log(data);
				$scope.url = "";
			})
			.error(function(data,status,headers,config) {
				// called asynchronously if an error occurs
				// or server returns response with an error status.
			});
	}
}]);
