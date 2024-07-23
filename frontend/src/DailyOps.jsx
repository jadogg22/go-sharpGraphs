import React, { useEffect } from 'react';
import PerformanceTable1 from './components/PerformanceTable1';
import PerformanceTable2 from './components/PerformanceTable2';

const DailyOps = () => {
  
    useEffect(() => {
        const fetchData = async () => {
        try {
            const response = await fetch('/api/Dispatch/Week_to_date/');
            const data = await response.json();
            console.log(data);
        } catch (error){
            console.error(error);
        }
    };

    fetchData();
}, []);

    {// rate per truck per day under 1000 rptpd 
        // only red if under 450
    }
    
return (
    <>
    <div className="p-6">
    <h1 className="text-2xl font-bold text-center p-6">Daily Operations Week to Data</h1>
    <PerformanceTable1 />
    <PerformanceTable2 />
   </div>

    </>
);
};
export default DailyOps;
