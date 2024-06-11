import React from 'react';

const PerformanceTable2 = () => {
  const textData = `OTR:   MPTD Goal: 565 RPTPD Goal: $1250.00   DH% Goal: <10% Service Goal: 98%+
LOCAL:   MPTD Goal: 184 RPTPD Goal: $500.00   DH% Goal: <45% Service Goal: 98%+
TEXAS:   MPTD Goal: 272 RPTPD Goal: $606.00   DH% Goal: <30% Service Goal: 98%+
Day 6/20 for month of June`;

  const tableData = [
    {
      driverManager: 'FELICIA MCMICHAEL',
      averageMPTPD: 539,
      averageRPTPD: 1304.93,
      deadheadPercentage: 6.8,
      orderOTP: 81.25,
      stopOTP: 88.03,
      averageMPTPDNeededToMakeGoal: 591,
    },
    {
      driverManager: 'TRINA SEPULVEDA',
      averageMPTPD: 407,
      averageRPTPD: 974.08,
      deadheadPercentage: 6.5,
      orderOTP: 67.86,
      stopOTP: 83.78,
      averageMPTPDNeededToMakeGoal: 723,
    },
    {
      driverManager: 'KEN MOSS JR',
      averageMPTPD: 518,
      averageRPTPD: 1231.54,
      deadheadPercentage: 11.2,
      orderOTP: 86.33,
      stopOTP: 93.85,
      averageMPTPDNeededToMakeGoal: 612,
    },
    {
      driverManager: 'LINDSAY WORKMAN',
      averageMPTPD: 492,
      averageRPTPD: 1238.08,
      deadheadPercentage: 7.3,
      orderOTP: 86.49,
      stopOTP: 93.81,
      averageMPTPDNeededToMakeGoal: 638,
    },
    {
      driverManager: 'ROCHELLE GENERA',
      averageMPTPD: 180,
      averageRPTPD: 680.88,
      deadheadPercentage: 38.2,
      orderOTP: 89.9,
      stopOTP: 94.13,
      averageMPTPDNeededToMakeGoal: 188,
    },
    {
      driverManager: 'STEPHANIE BINGHAM',
      averageMPTPD: 161,
      averageRPTPD: 400.45,
      deadheadPercentage: 23.3,
      orderOTP: 57.14,
      stopOTP: 71.43,
      averageMPTPDNeededToMakeGoal: 383,
    },
    {
      driverManager: 'TRACY RIGBY',
      averageMPTPD: 536,
      averageRPTPD: 1319.78,
      deadheadPercentage: 5.1,
      orderOTP: 86.67,
      stopOTP: 94.74,
      averageMPTPDNeededToMakeGoal: 594,
    },
  ];

  const getColorClass = (value, goal) => {
    if (value >= goal) {
      return 'bg-green-200';
    } else if (value >= goal * 0.9) {
      return 'bg-yellow-200';
    } else {
      return 'bg-red-200';
    }
  };

  return (
    <div>
      <pre className="bg-gray-200 p-4 rounded-md mb-4">{textData}</pre>
      <table className="table-auto w-full">
        <thead>
          <tr className="bg-gray-800 text-white">
            <th className="px-4 py-2">DRIVER MANAGER</th>
            <th className="px-4 py-2">Average MPTPD</th>
            <th className="px-4 py-2">Average RPTPD</th>
            <th className="px-4 py-2">DH%</th>
            <th className="px-4 py-2">ORDER OTP</th>
            <th className="px-4 py-2">STOP OTP</th>
            <th className="px-4 py-2">AVG MPTPD Needed to Make Goal</th>
          </tr>
        </thead>
        <tbody>
          {tableData.map((item, index) => (
            <tr key={index} className={index % 2 === 0 ? 'bg-gray-100' : 'bg-white'}>
              <td className="border px-4 py-2">{item.driverManager}</td>
              <td className={`border px-4 py-2 text-right ${getColorClass(item.averageMPTPD, 565)}`}>
                {item.averageMPTPD}
              </td>
              <td className={`border px-4 py-2 text-right ${getColorClass(item.averageRPTPD, 1250)}`}>
                {item.averageRPTPD.toFixed(2)}
              </td>
              <td className={`border px-4 py-2 text-right ${getColorClass(100 - item.deadheadPercentage, 90)}`}>
                {item.deadheadPercentage.toFixed(1)}%
              </td>
              <td className={`border px-4 py-2 text-right ${getColorClass(item.orderOTP, 98)}`}>
                {item.orderOTP.toFixed(2)}%
              </td>
              <td className={`border px-4 py-2 text-right ${getColorClass(item.stopOTP, 98)}`}>
                {item.stopOTP.toFixed(2)}%
              </td>
              <td className="border px-4 py-2 text-right">{item.averageMPTPDNeededToMakeGoal}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default PerformanceTable2;