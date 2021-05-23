var mTestApp = angular.module('mTestApp');
mTestApp.factory('toPdf', function () {
	return {
		saveAsHtml: function(data, sum,totalsum) {
		  	var rowArr = eval(data);
		  	var table2 = "<style> .result {box-shadow: 6px 6px 0px #c7c4c4;padding: 5px;overflow: hidden;border: 1px solid #000;margin: 0 auto;width: 1000px;} .zob {clear: both;background: #eeeded;} .proc {clear: both;}.pos {clear: both;border-bottom: 1px solid #000;} "
		  	var table3 = ".calc {border-bottom: 1px solid #999;width: 100%;clear: both; text-align: right;}.subname{width:80%;}.subnum{width:19%;text-align: right;}.subsum{background: #dcdcdc;width:100%;}.dia {clear: both;display: block;height: 65px!important;border: 1px solid #999;margin-left:3px;}.zobov:after {content: '.'; display: block; height: 0; clear: both; visibility: hidden;}"
		  	var table4 = ".etap{border-bottom: 2px solid #000;}.etap div{line-height: 22px;height: 21px;} .floatleft{float:left;} .floatright{float:right;}.totalsum {font-size: 1.2em;} .organ, .akt {border-bottom: 1px solid #000;} h4{margin:0px;}.zobov {height: auto!important;padding-left: 3px;clear: both;}</style>"
		  	var table =  "<div class='result'><h2 style='text-align:center;margin: 0px;'>Результат М-ТЕСТ:</h2><h4 class='organ'><b>Орган влади/організація: </b>"+rowArr[0]['subj']+"</h4><h4 class='akt'><b>Назва рег. акту: </b>"+rowArr[0]['regact']+"</h4>"
		  	var summa = '<div class="totalsum"><div class="floatleft"><b>Загальні витрати:</b></div><div class="floatright">(<b>Кі: </b>'+ rowArr[0]['ki']+')*(<b>Сума:</b> ' +(totalsum/rowArr[0]['ki']) + ' грн) = <b>'+totalsum+' грн</b></div></div>'
		  	for(var i = 0; i < rowArr.length; i++) { 
		    	table += '<div class="etap"><div class="zob"><b>Інф. зобов’язання:</b> ' + rowArr[i]['zob'] + '</div>';
		      	var col = rowArr[i]['columns'][0];
		        for(var j = 0; j < col.length; j++) {
		          	if (col[j]['type'] == 'itemplus'){
		            	table += '<div class="zobov"><div class="zob"><b>Комплексна дія:</b> '+col[j]['zob']+'</div>'
		            	zob = col[j]['columns'][0];
		            	for(var k = 0; k < zob.length; k++) {
		              		table += '<div class="dia">'
			              	table += '<div class="proc"><b>Дія:</b> '+zob[k]['proc']+'</div>'
			              	table += '<div class="pos"><b>Посада:</b> '+zob[k]['pos']+'</div>'
			              	table += '<div class="calc"><b>Чі:</b> '+zob[k]['chi']+' * <b>ВЧі: </b>'+zob[k]['vchi']+' * <b>КРІ:</b> '+zob[k]['kri']+' + <b>ПВРі:</b> '+zob[k]['pvri']+' = '
			              	table += (zob[k]['chi'] * zob[k]['vchi'] * zob[k]['kri']+ parseInt(zob[k]['pvri']))+ ' грн</div>'
			              	table += '</div>'
		            	}
		            	table += '</div>'
		          	} else {
		          		table += '<div class="dia">'
		          		table += '<div class="proc"><b>Дія:</b> '+col[j]['proc']+'</div>'
		          		table += '<div class="pos"><b>Посада:</b> '+col[j]['pos']+'</div>'
		          		table += '<div class="calc"><b>Чі:</b> '+col[j]['chi']+' * <b>ВЧі: </b>'+col[j]['vchi']+' * <b>КРІ:</b> '+col[j]['kri']+' + <b>ПВРі:</b> '+col[j]['pvri']+' = '
		         	 	table += (col[j]['chi'] * col[j]['vchi'] * col[j]['kri']+ parseInt(col[j]['pvri']))+ ' грн</div>'
		          		table += '</div>'
		        		}
		        	}
		      		table += '<div class="subsum"><div class="subname floatleft"><b>Витрати:</b> </div><div class="subnum floatleft"><b>'+sum[rowArr[i]['id']-1]+' грн</b> </div></div></div>'
		    	};
		  		return table2+table3+table4+table+summa+'</div>';
			},
		corruptEmail: function(data) {
			var fields = eval(data);
			output = '';
			for (var i =0; i < Object.keys(data).length; i++){
				var row = Object.keys(data)[i]
				if (row != 'drsu' && row != 'comment') {
					output += '<div border="2"><div class="row"><div class="col-sm-6">'+ fields[row]['question'] + '</div><div class="col-sm-2">' + fields[row]['ok'] + '</div><div class="col-sm-2">' + fields[row]['level'] + '</div><div class="col-sm-2">' + fields[row]['price'] + ' грн</div></div>'
					output += '<div class="row"><div class="col-sm-4">'+ fields[row]['norma1'] + '</div>' + '<div class="col-sm-4">'+ fields[row]['norma2'] + '</div>' + '<div class="col-sm-4">'+ fields[row]['norma3'] + '</div></div></div>'
				} else if (row == 'corrupSum') {
					'<div><div class="col-md-8">Сума:</div> <div class="col-md-4">' + fields[row] + '</div></div>'
				} else if (row == 'drsu') {
					output += '<div class="row" border="2"><div>Чи потрібне втручання ДРСУ?</div><div>' + fields[row] + '</div></div>'
				} else if (row == 'comment') {
					output += '<div class="row" border="2"><div class="col-sm-12">Коментар</div><div class="col-sm-12">' + fields[row] + '</div></div>'
				}
			}
			return output;

		}
		}
	})
