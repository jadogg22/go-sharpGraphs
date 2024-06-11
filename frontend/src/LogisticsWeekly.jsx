import React from 'react';
import LogisticsPerformaceTable from './components/LogisticsPerformaceTable';

const LogisticsWeekly = () => {
    return (
        <div className='p-8'>
            <h1 className="text-2xl font-bold  ">Logistics Weekly</h1>
            <LogisticsPerformaceTable />
        </div>
    );
};

export default LogisticsWeekly;