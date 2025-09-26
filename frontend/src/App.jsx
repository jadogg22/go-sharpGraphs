import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, NavLink } from 'react-router-dom';
import DailyOps from './DailyOps';
import LogisticsWeekly from './LogisticsWeekly';
import GraphYearlyRevenue from './GraphYearlyRevenue';

import EmptyMilesWrapper from './EmptyMilesWrapper';
import LogisticsMTD from './LogisticsMTD';
import RevPie from './RevPie';
import Dashboard from './Dashboard';
import OrdersPage from './OrdersPage';
import TransDrivers from './TransDrivers';
import LaneProfitability from './LaneProfitability';
import VacationExport from './VacationExport';

const App = () => {
  return (
    <Router>
      <div className="flex h-screen">
        <div className="w-64 bg-gray-800 text-white p-4 flex flex-col">
          <nav className="flex-grow flex flex-col justify-center">
            <ul>
              <li>
                <NavLink
                  to="/Dashboard"
                  className={({ isActive }) => isActive ? "block py-2 px-4 bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Dashboard
                </NavLink>
              </li>
              <li>
                <NavLink
                  end
                  to="/"
                  className={({ isActive }) => isActive ? "block py-2 px-4 bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Transporation Weekly Revenue
                </NavLink>
              </li>
              
              <li>
                <NavLink
                  to="/dailyOps"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Daily Ops
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/RevenuePie"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Revenue Pie
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/LogisticsMTD"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Logistics Month to Date
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/EmptyMiles"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Empty Miles
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/DriverManagers"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Drivers' Miles
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/logisticsWeekly"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Logistics Weekly
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/orders"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Sportsmans Invoice
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/lane-profitability"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Lane Profitability
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/vacation-export"
                  className={({ isActive }) => isActive ? "bg-gray-700" : "block py-2 px-4 hover:bg-gray-700"}
                >
                  Vacation Export
                </NavLink>
              </li>
            </ul>
          </nav>
        </div>
        <div className="flex-1 p-4 flex flex-col justify-center overflow-auto">
          <Routes>
            <Route path="/" element={<GraphYearlyRevenue />} />
            <Route path="/Dashboard" element={<Dashboard />} />
            
            <Route path="/dailyOps" element={<DailyOps />} />
            <Route path="/RevenuePie" element={<RevPie />} />
            <Route path="/LogisticsMTD" element={<LogisticsMTD />} />
            <Route path="/EmptyMiles" element={<EmptyMilesWrapper />} />
            <Route path="/DriverManagers" element={<TransDrivers />} />
            <Route path="/Orders" element={<OrdersPage />} />
            <Route path="/lane-profitability" element={<LaneProfitability />} />
            <Route path="/vacation-export" element={<VacationExport />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

const Home = () => (
  <div>
    <h1 className="text-2xl font-bold">Home Page</h1>
    <p>Welcome to the home page!</p>
  </div>
);

const About = () => (
  <div>
    <h1 className="text-2xl font-bold">About Page</h1>
    <p>Learn more about us on the about page.</p>
  </div>
);


export default App;
