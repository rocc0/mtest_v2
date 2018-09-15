var interceptor = function ($q, $location) {
    return {
        request: function (config) {//req
            if (!localStorage.getItem('token')) {
                localStorage.setItem('token', 'unknown');
            }
            return config;
        },

        response: function (result) {//res
            return result;
        },

        responseError: function (rejection) {
            if (rejection.status === 401 && rejection.config.method === 'POST') {
                $location.url('/u/login');
            }
            return $q.reject(rejection);
        }
    }
};

mTestApp.config(function($compileProvider, $interpolateProvider, $routeProvider, $locationProvider,$httpProvider) {
    $httpProvider.interceptors.push(interceptor);
    $compileProvider.aHrefSanitizationWhitelist(/^\s*(https?|ftp|mailto|tel|file|blob):/);
    $locationProvider.html5Mode(true);
    $interpolateProvider.startSymbol('{[{').endSymbol('}]}');
    $routeProvider
        .when('/', {
            title: 'Калькулятор',
            templateUrl: 'static/html/dnd_anon.html',
            controller: 'mTestController'
        })
        .when('/search', {
            title: 'Пошук АРВ',
            templateUrl: 'static/html/search.html',
            controller: 'searchController'
        })
        .when('/track/:mtest_id', {
            templateUrl: 'static/html/dnd_db.html',
            controller: 'mTestDBController'
        })
        .when('/u/register', {
            title: 'Реєстрація',
            templateUrl: '/static/html/auth/auth.register.view.html',
            controller: 'authRegisterController'
        })
        .when('/u/login', {
            title: 'Вхід',
            templateUrl: '/static/html/auth/auth.login.view.html',
            controller: 'authLoginController'
        })
        .when('/u/cabinet', {
            title: 'Кабінет користувача',
            templateUrl: '/static/html/auth/auth.cabinet.view.html',
            controller: 'userCabinetController',
            resolve: {
                mTestApp: function ($q) {
                    var defer = $q.defer();
                    defer.resolve();
                    return defer.promise;
                }
            }
        })
        .when('/u/activate/:hash', {
            title: 'Активація акаунту',
            templateUrl: '/static/html/auth/auth.activate.view.html',
            controller: 'authActivateController'
        })
        .when('/u/reset/:hash', {
            title: 'Відновлення доступу',
            templateUrl: '/static/html/auth/auth.reset.view.html',
            controller: 'authResetController'
        })
        .otherwise({redirectTo:'/'});

});

mTestApp.run(['$rootScope', '$route', function($rootScope, $route) {
    $rootScope.$on('$routeChangeSuccess', function(event, currentRoute, previousRoute) {
        document.title = "M-TECT | " + currentRoute.title;
    });
}]);