mTestApp.service('authService', function ($http) {
    /*jshint validthis: true */
    const baseURL = 'http://localhost:8888/';

    this.login = function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/login',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.logout = function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/logout',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.register =  function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/register',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.ensureAuthenticated = function(token) {
        return $http({
            method: 'GET',
            url: baseURL + 'api/v.1/u/cabinet',
            headers: {
                'Content-Type': 'application/json',
                Authorization: 'Bearer ' + token
            }
        });
    };
});

mTestApp.service('mtCrud', function ($http) {
    const baseURL = 'http://localhost:8888';

    this.readMtestFromDB = function (id) {
        return $http({
            method: 'GET',
            url: '/api/v.1/m/get/'+id,
        });
    }

    this.addMtest = function(newmtest, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/create",
            data: {
                name: newmtest.name,
                region: parseInt(newmtest.region.id),
                government: parseInt(newmtest.government.id)
            },
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.removeMtestItem = function (id, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/delete",
            data: {id:id},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };

    this.updateMtestItem = function (item, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/update",
            data: {mid:item.id, name:item.name,
                region:parseInt(item.region),govern: parseInt(item.govern)},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.updateMtestCalculations= function (id, item, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/update",
            data: {id:id,calculations: item},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
});