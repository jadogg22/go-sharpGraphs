# go-sharpGraphs

Here is my take on of the backend for my sharpGraphs website showing internal company revenue charts for better understanding towards trends of the market.

Right now the biggest problem is actually getting the attachments, I've figured out that with onedrive that I can automatically download files from my email. using powerautomate I Tried pretty hard to get the attachments in GO but there isn't very much documeneation for that yet so I went with the power automate method, 

so now that I can put that on my local system, There is now a second and third part that I need to develop.

## number 2

So now that we are successfully getting the files from our email, they are still stuck the harddrive of my windows pc... kinda unfortuate that it only works with windows but it is what it is. Basically Now I need to write simple little program that starts everyday at like 8:15 and parses the new file, then probably deletes it once its done. This should be simple enough but we'll see.

- This is finished, the only way that I was able to get it working was using powerautomate from microsoft that automatically puts it into a special folder that I have on my computer and this program that sits in the same directory and will parse and send the data to a post endpoint on the server. So its basically all automated.


## ammendments

So as i've been implementing the year_by_year data Ive realized that I do infact need to do send a slice of hashmaps that is going to be the best method for includeing or not including the years that are not valad like a future date.

another point is that Ive really been thinking about the best method for pulling from the db. I ended up running some benchmarks, I do a complicated sql command that gets all the rows in their proper order but kinda breaks and throws some errors because of weird weeks. Then I created a bach and transaction system that will string together a query for every year. Its about 3 times slower but there are no errors and it is much more configureable. So I think I'm gonna just go with that method and as long as the query doesnt return an err we can add this month to the slice 

- I got around this with that additional table, I just select star and for every record I get the year and the week, go to that specific week and use the year as a key and just put in the data. Then for the data thats not in the table yet I just figure out the gap where the data doesn't exist query the other table for the missing data by week and if it exists I add it! 


## Coded Revenue

In previous verisons of Coded revenue, I just sent all the data and then I had javascript sort out all the data, compressing it and compiling all the small ones. I figured that was a poor Idea so now I do it all on the backend and to make it even easier on myself I can add other things in the json like the Total Revenue and whenever I need that value... I just ask for it! long gone are the days where i'm just sending an array over json.

I actually really like this graph and I havn't nearly given it enough love. This is going to be my next area of updating quite a bit. So Ive already included the endpoint to exept perameters,

Where I'm going with this is that I would like to add some drop downs for months and years to select spacific timeframes. Maybe in the future Id like to be able to take some of the codes out as well from the data? 

## Dispachers

So The Week to Date in the dispachers was really probably one of the main reasons I even started this project. Luke the first thing that he does in the morning when he gets to the office is download this excel file open it up in powerbi and then run some numbers and send off a report to the higherups. I wanted to automate this proccess for him! 

Right now I have it pretty much working. Ive made it so the f:ile automatically gets downloaded to my computer when I receive it in my email from mclaud I also created a post endpoint to send the data to so I parse the data on my computer and send it to the server then once I do that it will automatically show up on the site. This way the report is automatically done and you can compair all the dispachers against one another.
