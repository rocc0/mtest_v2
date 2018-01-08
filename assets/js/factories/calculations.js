var mTestApp = angular.module('mTestApp');
mTestApp.factory('MathData', function () {
    return {
      mathSum: function(data,id) {
        var total = {};
        var subtotal = 0;
        var arr = eval(data);
        for(var i = 0; i < arr.length; i++) { 
          var itemid = arr[i]['id'];
          var col = arr[i]['columns'][0];
            for(var j = 0; j < col.length; j++) {
              if (col[j]['type'] == 'itemplus'){
                var itemplus = col[j]['columns'][0];
                for(var k = 0; k < itemplus.length; k++){
                  if (itemplus[k]['subsum'] > 0){
                    subtotal += Math.round(itemplus[k]['subsum'])
                  } else {
                  subtotal += 0
                  }
                }
              } else {
                if (col[j]['subsum'] > 0){
                subtotal += col[j]['subsum']
              } else {
                subtotal += 0
                }
              }
            }
            total[itemid] = subtotal;
            subtotal *= 0;
          }
        return total
      },
      totalSum: function(data, func){

          if (data > 0){
            data = data
          } else {
            data = 1
          }
          var val = func();
          var valKeys = Object.keys(func());
          var total = 0;
          for (i =0; i < valKeys.length; i++){
            total += parseInt(val[valKeys[i]])
          }
          return total * data

      },
      awgSum: function(data, ki){

          if (ki > 0){
            ki = ki
          } else {
            ki = 1
          }
          var total = 0
          for (i =0; i < data.length; i++){
            total += parseInt(data[i]['awgsub'])
          }
          return total * ki

      },
      getNum: function(data){
        var arr = eval(data['1']);
        var max = 0
        for (i =0; i < arr.length; i++){
          if (parseInt(arr[i]['id']) > max){
            max = parseInt(arr[i]['id'])
          }
        }
        return max
      }
    }
});