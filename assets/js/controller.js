var mTestApp = angular.module("mTestApp", ["ngAnimate", "ngSanitize", "ngCookies", "ngRoute", "ngTagsInput",
    "dndLists", "ui.bootstrap", "ngAnimate", "sticky", "switcher"]);

//----------------------------------------------------------------------------------------------------------
//---------------------------------------------CONTROLLER WITH LOCAL ----------------------------------------
//----------------------------------------------------------------------------------------------------------

mTestApp.controller("mTestController", function ($scope, $sce, $http, $cookies, $window,
                                                 $location, LS, SubmitData, MathData, ModalWin, $interval) {
    $http.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
    //----------------------------------------------------------------------------------------------------------
    //---------------------------------------------DATA---------------------------------------------------------
    //----------------------------------------------------------------------------------------------------------

    this.value = LS.getData();
    this.data = angular.fromJson(this.value);
    $scope.disableSticking = false;
    this.num = MathData.getNum(this.data);

    $scope.allowed = {dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item']};
    $scope.models = {
        selected: null,
        templates: [{
            type: "container",
            id: this.num + 1,
            columns: [
                []
            ],
            name: "Додати інф. вимогу"
        }, {
            type: "itemplus",
            id: 3,
            columns: [
                []
            ],
            name: "Додати складову інф. вимоги"
        }, {
            type: "item",
            id: 3,
            name: "Додати дію"
        }],
        dropzones: this.data
    };

    setInterval(function () {
        var lsData = angular.toJson($scope.models.dropzones);
        LS.setData(lsData);
    }, 5000);

    $http({
        method: 'GET',
        url: '/api/v.1/adm_actions'
    }).then(function (response) {
        $scope.adm_actions = response.data.actions;
    });


    $http({
        method: 'GET',
        url: '/static/questions.json'
    }).then(function (response) {
        $scope.questions = response.data;
    }).catch(function (reason) {
        console.log(reason)
    });


    $scope.isLoggedIn = 234;

    $scope.mdata = function () {
        return angular.toJson($scope.models.dropzones);
    };

    $scope.effectPopover = $sce.trustAsHtml('Критично – реалізація корупційного ризику призведе до зупинки діяльності <br> ' +
        'Суттєво - реалізація корупційного ризику призведе до суттєвих витрат грошей та/або часу співставних з розміром щомісячного прибутку');

    var trusted = {};

    $scope.getPopoverContent = function (content) {
        return trusted[content] || (trusted[content] = $sce.trustAsHtml(content));
    };

    //--------------------------------------------- ADD-REMOVE ROW ---------------------------------------------
    $scope.addRow = function (item, text) {
        var quest = {
            text: text,
            "yesno": [{"name": "Так", "value": "yes"}, {"name": "Ні", "value": "no"}],
            "remove": "block"
        };
        item.push(quest);
    };

    $scope.removeRow = function (item, index) {
        item.splice(index, 1);
    };

    $scope.deleteTrash = function (where, stndrt, real) {
        return where.slice(stndrt, real)
    };
    //---------------------------------------------SUM----------------------------------------------------------

    $scope.Sum = function (id) {
        return MathData.mathSum($scope.models.dropzones['1'], id)
    };

    $scope.totalSum = function () {
        return MathData.totalSum($scope.models.dropzones['1'][0]['ki'], this.Sum)
    };

    $scope.corr_sum = function (data) {
        var total = 0;
        for (var k in data) {
            total += Number(data[k])
        }
        return total;
    };

    $scope.addToSlider = function (res, eft) {
        if (res == 1 && eft == 1) {
            return 1
        } else if (res == 1 && eft > 2) {
            return 100
        } else {
            return 0
        }

    };

    $scope.getLength = function (obj) {
        if (obj != null) {
            return obj.length
        } else {
            setTimeout($scope.getLength, 300);
        }

    };
    //---------------------------------------------MODAL--------------------------------------------------------

    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };



    $scope.hideQuestion = function (dep, val) {
        if (dep === 1 && val === 0) {
            return "none"
        } else if (dep === 1 && val === 1) {
            return "block"
        } else {
            return "block"
        }
    };
    $scope.payContent = $sce.trustAsHtml('Оплата – зазвичай середня  вартість людино-години роботи персоналу, який\
     виконує дії необхідні для  виконання ІВ.<br>\
  За бажанням можна використовувати  заробітну плату однієї години відповідного Спеціаліста або Керівника. <br>\
  Для ФОПів  це може бути  розрахунок вартості  людино-години на основі річного обсягу реалізації сектору (детальніше див.  Посібник з М-Тесту).<br>\
  В віконце  оплата вводиться  за місяць і перераховується автоматично в оплату за годину.');


    //---------------------------------------------CORR_BAR--------------------------------------------------------
    $scope.hideCritical = function (val, length) {
        if (val >= 100 || val > length / 2) {
            return "block"
        } else {
            return "none"
        }
    };

    $scope.corr_bar = function (total, length) {
        var a = '';
        if (total > length / 2) {
            a = 'danger';
            $scope.barText = "Критична корупційна складова!";
        } else {
            $scope.barText = "Кількість корупційних ризиків:" + total;
            a = 'success'
        }
        return a
    };

    $scope.valueToText = function (val, type) {
        if (type === "yn") {
            if (val === 1 || val === "Yes") {
                return "Так"
            } else if (val === 0 || val === "No") {
                return "Ні"
            } else if (val === "idk") {
                return "Незнаю"
            } else {
                return "Не заповнено!"
            }
        } else {
            if (val === 1) {
                return "Суттєво"
            } else if (val === 0) {
                return "Несуттєво"
            } else if (val === 100) {
                return "Критично"
            } else {
                return "Не заповнено!"
            }
        }
    };

    $interval(function() {
        var a = document.getElementById('report').innerHTML;
        var blob = new Blob(["\ufeff", a], {
            type: 'text/html'
        });
        $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
    }, 1000);

    //---------------------------------------------RESET--------------------------------------------------------
    $scope.reset = function () {
        $scope.models.dropzones = angular.copy({
            "1": [{
                "type": "container", "id": 3, "columns": [[{
                    "type": "itemplus", "id": 3,
                    "columns": [[{"type": "item", "id": 3, "name": "Додати дію", "subsum": 0}, {
                        "type": "item",
                        "id": 6,
                        "name": "Додати дію",
                        "subsum": 0
                    }]],
                    "name": "Додати складову інф. вимоги"
                }]], "name": "Додати інф. вимогу", "contsub": 0
            },
                {
                    "type": "container", "id": null, "columns": [[{
                        "type": "itemplus",
                        "id": 4,
                        "columns": [[{"type": "item", "id": 3, "name": "Додати дію", "subsum": 0},
                            {"type": "item", "id": 4, "name": "Додати дію", "subsum": 0}]],
                        "name": "Додати складову інф. вимоги"
                    }]], "name": "Додати інф. вимогу", "contsub": 0
                }]
        });
    }

});

mTestApp.controller("mTestDBController",
    function ($scope, $timeout, $sce, $http, $cookies, $location, $routeParams, DB_data,
              $window, LS, mtCrud, MathData, ModalWin, $interval) {
        const token = localStorage.getItem('token');
        $scope.allowed = {dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item']};

        $scope.params = $routeParams;

        var trusted = {};
        $scope.getPopoverContent = function (content) {
            return trusted[content] || (trusted[content] = $sce.trustAsHtml(content));
        };
        //----------------------------------------------------------------------------------------------------------
        //---------------------------------------------DATA---------------------------------------------------------
        //----------------------------------------------------------------------------------------------------------
        mtCrud.readMtestFromDB($scope.params.mtest_id)
            .then(function (response) {
                $scope.mtestData = response.data.mtest;
                $scope.executors = angular.fromJson(response.data.mtest.executors);
                this.val = angular.fromJson(response.data.mtest.calculations);
                $scope.modelsdb = {
                    selected: null,
                    templates: [{
                        type: "container",
                        id: this.num + 1,
                        columns: [
                            []
                        ],
                        name: "Додати інф. вимогу"
                    }, {
                        type: "itemplus",
                        id: 3,
                        columns: [
                            []
                        ],
                        name: "Додати складову інф. вимоги"
                    }, {
                        type: "item",
                        id: 3,
                        name: "Додати дію"
                    }],
                    dropzones: this.val
                };
            });

        //-------------------------------------------- FOR TYPEHEAD ------------------------------------------------

        $http({
            method: 'GET',
            url: '/api/v.1/adm_actions'
        }).then(function (response) {
            $scope.adm_actions = response.data.actions;
        });

        $http({
            method: 'GET',
            url: '/static/questions.json'
        }).then(function (response) {
            $scope.questions = response.data;
        }).catch(function (reason) {
            console.log(reason)
        });


        //-------------------------------------------- GROUP MATH ------------------------------------------------
        $scope.doGroupMath = function () {
            var execs = $scope.executors;
            var infs = $scope.modelsdb.dropzones['1'];

            for (var devs_get = []; devs_get.length < infs.length; devs_get.push([]));
            for (var p in execs) {
                let dev_data;
                console.log(execs[p].checked);
                if (execs[p].checked === true) {
                    $http({method: 'GET', url: '/api/v.1/m/get/' + execs[p].mid })
                        .then(function (response) {
                            dev_data = JSON.parse(response.data.mtest.calculations)
                        });
                    setTimeout(function () {
                        for (var j = 0; j < dev_data['1'].length; j++) {
                            if (dev_data['1'][j]['contsub'] !== 0) {
                                devs_get[j].push(dev_data['1'][j]['contsub'])
                            }
                        }
                    }, 50)
                }
            }
            $scope.devs_get = devs_get
        };

        setTimeout(function () {
            $scope.doGroupMath()
        }, 150);


        setTimeout(function () {
            $scope.sumDevsInfs = function (x) {
                var tmp_sum = 0;
                for (var i = 0; i < x.length; i++) {
                    tmp_sum += x[i]
                }
                return tmp_sum / x.length
            }
        }, 200);


        //-------------------------------------------- END GROUP MATH ------------------------------------------------

        //-------------------------------corruption part

        $scope.corr_sum = function (data) {
            var total = 0;
            for (var k in data) {
                total += Number(data[k])
            }
            return total;
        };


        $scope.corr_bar = function (total, length) {
            var a = '';
            if (total > length / 2) {
                a = 'danger';
                $scope.barText = "Критична корупційна складова!";
            } else {
                $scope.barText = "Кількість корупційних ризиків:" + total;
                a = 'success'
            }
            return a
        };

        $scope.addRow = function (item, text) {
            var quest = {
                text: text,
                "yesno": [{"name": "Так", "value": "yes"}, {"name": "Ні", "value": "no"}],
                "remove": "block"
            };
            item.push(quest);
        };

        $scope.hideCritical = function (val, length) {
            if (val >= 100 || val > length / 2) {
                return "block"
            } else {
                return "none"
            }
        };

        $scope.removeRow = function (item, index) {
            item.splice(index, 1);
        };
        $scope.addToSlider = function (res, eft) {
            if (res === 1 && eft === 1) {
                return 1
            } else if (res === 1 && eft > 2) {
                return 100
            } else {
                return 0
            }

        };

        $scope.hideQuestion = function (dep, val) {
            if (dep === 1 && val === 0) {
                return "none"
            } else if (dep === 1 && val === 1) {
                return "block"
            } else {
                return "block"
            }
        };

        $scope.valueToText = function (val, type) {
            if (type === "yn") {
                if (val === 1 || val === "Yes") {
                    return "Так"
                } else if (val === 0 || val === "No") {
                    return "Ні"
                } else if (val === "idk") {
                    return "Незнаю"
                } else {
                    return "Не заповнено!"
                }
            } else {
                if (val === 1) {
                    return "Суттєво"
                } else if (val === 0) {
                    return "Несуттєво"
                } else if (val === 100) {
                    return "Критично"
                } else {
                    return "Не заповнено!"
                }
            }
        };

        //end corruption part --------------------

        //----------------------------------------------------------------------------------------------------------
        //---------------------------------------------SAVE TO DB --------------------------------------------------
        //----------------------------------------------------------------------------------------------------------

        $scope.mdata = function () {
            return angular.toJson($scope.modelsdb.dropzones);
        };
        $scope.saveMtestToDB = function () {
            mtCrud.updateMtestCalculations($scope.params.mtest_id, $scope.mdata(), token)
                .then(function (value) {
                    console.log(value.data)
                })
        };
        $scope.updateExecutors = function () {
            var executors = angular.toJson($scope.executors);
            mtCrud.updateMtestExecutors($scope.params.mtest_id, executors, token)
                .then(function (value) {
                    $scope.doGroupMath()
                    console.log(value.data)
                }).catch(function (err) {  console.log(err) })
        };
        //---------------------------------------------SUM----------------------------------------------------------

        $scope.SumDb = function (id) {
            return MathData.mathSum($scope.modelsdb.dropzones['1'], id)
        };
        $scope.totalSumDb = function () {
            return MathData.totalSum($scope.modelsdb.dropzones['1'][0]['ki'], $scope.SumDb)
        };

        $scope.awgSumDb = function () {
            return MathData.awgSum($scope.modelsdb.dropzones['1'], $scope.modelsdb.dropzones['1'][0]['ki'])
        };

        //---------------------------------------------MODAL--------------------------------------------------------
        $scope.openModal = function (size, template, controller) {
            return ModalWin.openModal(size, template, controller)
        };


        //---------------------------------------------alert------------------------------------------------------


        $scope.getClass = 'not_disabled';
        $scope.getClassTwo = 'disabled';

        $scope.setClasses = function () {
            $scope.getClass = 'disabled';
            $scope.getClassTwo = 'not_disabled';
            $timeout(function () {
                $scope.getClass = 'not_disabled';
                $scope.getClassTwo = 'disabled';
            }, 1000)
        };


        //---------------------------------------------PDF--------------------------------------------------------
        $interval(function () {
            var a = document.getElementById('report').innerHTML;
            var blob = new Blob(["\ufeff", a], {
                type: 'text/html'
            });
            $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
        }, 1000);


//---------------------------------------------reset------------------------------------------------------
        $scope.resetdb = function () {
            mtCrud.readMtestFromDB($scope.params.mtest_id).then(function (response) {
                var val = angular.fromJson(response.data.mtest.calculations);
                $scope.modelsdb.dropzones = angular.copy(
                    val
                )
            })
        };

        $scope.effectPopover = $sce.trustAsHtml('Критично – реалізація корупційного ризику призведе до зупинки діяльності <br> ' +
            'Суттєво - реалізація корупційного ризику призведе до суттєвих витрат грошей та/або часу співставних з розміром щомісячного прибутку');
        $scope.payContent = $sce.trustAsHtml('Оплата – зазвичай середня  вартість людино-години роботи персоналу, який\
     виконує дії необхідні для  виконання ІВ.<br>\
  За бажанням можна використовувати  заробітну плату однієї години відповідного Спеціаліста або Керівника. <br>\
  Для ФОПів  це може бути  розрахунок вартості  людино-години на основі річного обсягу реалізації сектору (детальніше див.  Посібник з М-Тесту).<br>\
  В віконце  оплата вводиться  за місяць і перераховується автоматично в оплату за годину.');
    });

mTestApp.controller("authLoginController", function ($scope, $timeout, $location, authService) {
    $scope.template_name = true;
    $scope.user = {};
    $scope.onLogin = function () {
        authService.login($scope.user)
            .then(function (user) {
                localStorage.setItem('token', user.data.token);
                $location.path('/u/cabinet');
            })
            .catch(function (err) {
                console.log(err);
                $scope.message = "Невірний логін або пароль, спробуйте ще раз";
                $timeout(function () {
                    $scope.message = ""
                }, 2000);
            });
    };
    $scope.sendResetRequest = function () {
        authService.resetpass($scope.resetpass.email)
            .then(function () {
                $scope.message = "На ваш Email відправлено посилання для відновлення паролю"
            })
            .catch(function () {
                $scope.message = "Даний Email не зареєстровано";
                $timeout(function () {
                    $scope.message = ""
                }, 2000);
            })
    }
});

mTestApp.controller("userCabinetController", function ($scope, $http, $location, $rootScope,
                                                       $timeout, authService, mtCrud, ModalWin) {
    $scope.changepass = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token)
            .then(function (user) {
                if (user.status === 200) {
                    $scope.userdata = user.data.data;
                    $scope.records = user.data.data.records;
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function (err) {
                console.log(err);
                $location.path('/u/login');
            });
    }



    //load governments and regions
    $http({
        method: 'GET',
        url: '/api/v.1/regions',
    }).then(function (response) {
        $scope.regions = response.data.regions
    }).catch( function (reason) {
        console.log(reason)
        });

    $http({
        method: 'GET',
        url: '/api/v.1/govs',
    }).then(function (response) {
        $scope.governs = response.data.govs
    }).catch( function (reason) {
        console.log(reason)
    });
    //end load governments and regions

    //format label for typehead on select
    $scope.formatLabel = function(model, index, itmtype) {
        console.log(model, index);
        for (var i=0; i< $scope[itmtype+'s'].length; i++) {
            if (model.id === $scope[itmtype+'s'][i].id) {
                $scope.records[index][itmtype] = $scope[itmtype+'s'][i].id
            }
        }
    };
    //end format label for typehead on select

    $scope.changeUserField = function (field, id, value) {
        console.log(field, id, value);
        $http({
            method: 'POST',
            url: "/api/v.1/u/edituser",
            data: {field: field, data: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
    //--------------CRUD mtest------------------------
    //add executor to mtest item
    $scope.addExecutor = function (title, email, region, gov, dev_mid) {
        mtCrud.addMtestExecutor(title, email,region, gov, dev_mid, token).then(function (response) {
            $scope.records[dev_mid].executors[email] = {email: email, mid: response.data.mid}
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeExecutor = function (email, exIndex, devIndex) {
        mtCrud.removeMtestExecutor(email, exIndex, devIndex, token).then(function () {
            delete $scope.records[devIndex].executors[exIndex]
        }).catch(function (err) {
            console.log(err)
        });
    };
    //end add executor to mtest item
    $scope.addMtest = function (newmtest) {
        mtCrud.addMtest(newmtest, token).then(function (response) {
            $scope.records = response.data.records;
            $scope.len_records = angular.toJson($scope.records).length
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeMtestItem = function (id) {
        console.log(id);
        mtCrud.removeMtestItem(id, token)
            .then(function () {
                delete $scope.records[id];
                $scope.len_records = angular.toJson($scope.records).length
            }).catch(function (err) {
            console.log(err)
        });
    };

    $scope.updateMtestItem = function (item) {
        mtCrud.updateMtestItem(item, token)
            .then(function (response) {
            }).catch(function (err) {
            console.log(err)
        });
    };

    // popover for setting, mail, add mtest etc..
    $scope.dynamicPopover = {
        isOpen: {},
        title: "heh",

        templateUrl: {},

        open: function open(index, template, tp) {
            $scope.dynamicPopover.templateUrl[index + tp] = template;
            $scope.dynamicPopover.isOpen[index + tp] = true;
            $scope.dynamicPopover.data = 'Hello!';
        },

        close: function close(index, tp) {
            $scope.dynamicPopover.isOpen[index + tp] = false;
        }
    };
    // end popover for setting, mail, add mtest etc..

    //count len of records array for showing-hiding empty array banner
    var countUp = function () {
        $scope.len_records = angular.toJson($scope.records).length
    };
    $timeout(countUp, 300);
    //end count len of records array for showing-hiding empty array banner
});

mTestApp.controller("authRegisterController", function ($scope, authService,$location,$timeout) {

    $scope.user = {
        password: "",
        confirmPassword: ""
    };
    $scope.onRegister = function () {
        authService.register($scope.user)
            .then(function (response) {
                $location.path('/u/login');
            })
            .catch(function (err) {
                $scope.show_err = 1;
                $scope.err_status = err.data.title;
                $timeout(function () {
                    $scope.show_err = 0
                }, 2000);
                console.log(err);
            });
    };
});

mTestApp.controller("menuController", function ($scope, $rootScope,$location, authService, ModalWin) {
    $rootScope.isLoggedIn = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token)
            .then(function (user) {
                if (user.status === 200) {
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function () {
                $rootScope.isLoggedIn = false;
            });
    }
    $scope.onLogout = function () {
        localStorage.removeItem('token');
        $rootScope.isLoggedIn = false;
        $location.path('/u/login');
    };
    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
});

mTestApp.controller("searchController", function ($scope, $http) {
    $scope.currentPage = 0;
    //elastic search
    $scope.query = {
        "from": 0, "size": 10,
        "query": {
            "bool": {
                "should": {
                    "multi_match": {
                        "query": $scope.phrase,
                        "fields": ["name"]
                    }
                }
            }
        }
    };
    $scope.addPhrase = function () {
        $scope.query.query.bool.should.multi_match.query = $scope.phrase
    };

    $scope.addTerm = function (field, data) {
        if (!$scope.query.query.bool.filter) {
            $scope.query.query.bool.filter = {
                "bool": {
                    "must": []
                }
            }
        }
        var obj = {};
        var arr = $scope.query.query.bool.filter.bool.must;
        obj[field] = data;
        if (arr.length === 0) {
            arr.push({"term": obj})
        } else {
            for (var i = 0; i < arr.length; i++) {
                if (arr[i].term.hasOwnProperty(field)) {
                    arr[i].term[field] = data;
                    break
                } else if (i === arr.length - 1 && !arr[i].term.hasOwnProperty(field)) {
                    arr.push({"term": obj})
                }
            }
        }
    };
    $scope.doSearch = function () {
        $http({
            method: 'POST',
            url: "http://mtest.com.ua:9200/mtests/_search",
            data: $scope.query
        }).then(function (response) {
            $scope.results = response.data;
        }).catch(function (err) {
            console.log(err)
        })
    };

    //load governments and regions
    $http({
        method: 'GET',
        url: '/api/v.1/regions',
    }).then(function (response) {
        $scope.regions = response.data.regions
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/govs',
    }).then(function (response) {
        $scope.governs = response.data.govs
    }).catch( function (reason) {
        console.log(reason)
    });
    //end load governments and regions

});

mTestApp.controller("authActivateController", function ($scope, $routeParams,$http) {
    const baseURL = 'http://mtest.com.ua/';
    $http({
        method: 'GET',
        url: baseURL + 'api/v.1/u/activate/' + $routeParams.hash ,
    }).then(function (response) {
        $scope.message = "Акаунт успішно активовано"
    }).catch(function (err) {
        $scope.err_message = "Посилання не існує"
    });

});

mTestApp.controller("authResetController", function ($scope, $routeParams, $http, $location, authService) {
    const baseURL = 'http://mtest.com.ua/';
    var hash = $routeParams.hash;
    $scope.user = {};
    authService.checkhash(hash)
        .then(function (value) {
        })
        .catch(function (reason) {
            console.log(reason);
            $location.url('/');
        });
    $scope.onReset = function () {
        authService.setnewpass($scope.user.password, hash)
            .then(function () {
                $scope.message = "Пароль успішно змінено"
            }).catch(function () {
            $scope.err_message = "Помилка: посилання не існує"
        });
    }

});