
async function getMedian(){
	$('#median').text("Loading...");
	const jsonMedian= await callRestAPI(`/getMedian` );
	const median = JSON.parse(jsonMedian);
	console.log(jsonMedian)
	$('#median').text("");
	$('#median').append("<li> Median = " + median.Value + " </li>")
	$('#median').append("<li> SumList: " + median.SumList + " </li>");
}

var pageIndex = 1
async function get10OldestSum(){
		$('#sum10Oldest').text("Loading...");
		const jsonSum = await callRestAPI(`/get10Oldest/${pageIndex}` );
		const sum = JSON.parse(jsonSum);
		console.log(jsonSum)
		$('#sum10Oldest').text("");
		$('#sum10Oldest').append("<li>Sum = " + sum.Value + " </li>");
		pageIndex++;
}

async function updateLatest20Every5Second(){
	$('#20latest').text("Loading...");
	var interval = setInterval(async () => {		
		const pairs = await callRestAPI("/get20Latest");
		const pairsList = JSON.parse(pairs);
		console.log(pairs)
		$('#20latest').text("");
		if(pairsList != null){
			pairsList.forEach(function (pair) {
				try {
					$('#20latest').append("<li>" + pair.Value +" </li>"); 
				} catch (e) {
					console.error(e)
				}
			})
		}
		else{
			clearInterval(interval)
		}
	},5000);
	
}

async function countDownFrom(max){
	var count = max;
	var interval = setInterval(async () => {
		if(count > 0) {
			console.log(count)
			count--;
		}	else {
			count = max
		}
	}
	,1000)
}

function callRestAPI(url, method = "GET") {
  return new Promise(function(resolve, reject) {
    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(e) {
      if (xhr.readyState === 4) {
        if (xhr.status === 200) {
          resolve(xhr.response)
        } else {
          reject(xhr.status)
        }
      }
    }
    xhr.ontimeout = function () {
      reject('timeout')
    }
    xhr.open(method, url, true)
    xhr.send()
  })
}
  

