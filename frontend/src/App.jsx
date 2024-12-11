import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, NavLink } from 'react-router-dom';
import DailyOps from './DailyOps';
import LogisticsWeekly from './LogisticsWeekly';
import GraphYearlyRevenue from './GraphYearlyRevenue';
import LogYearlyRevenue from './LogYearlyRevenue.jsx';
import EmptyMilesWrapper from './EmptyMilesWrapper';
import LogisticsMTD from './LogisticsMTD';
import RevPie from './RevPie';
import Dashboard from './Dashboard';
import OrdersPage from './OrdersPage';

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
                  activeClassName="p-4 bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Dashboard
                </NavLink>
              </li>
              <li>
                <NavLink
                  end
                  to="/"
                  activeClassName="p-4 bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Transporation Weekly Revenue
                </NavLink>
              </li>
              <li>
                <NavLink
                  end
                  to="/LogisticsYearly"
                  activeClassName="p-4 bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Logistics Weekly Revenue
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/dailyOps"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Daily Ops
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/RevenuePie"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Revenue Pie
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/LogisticsMTD"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Logistics Month to Date
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/EmptyMiles"
                  activeClassName="bg-gray-700"
                  className="block py-2 px-4 hover:bg-gray-700"
                >
                  Empty Miles
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
            </ul>
          </nav>
        </div>
        <div className="flex-1 p-4 flex flex-col justify-center overflow-auto">
          <Routes>
            <Route path="/" element={<GraphYearlyRevenue />} />
            <Route path="/Dashboard" element={<Dashboard />} />
            <Route path="/LogisticsYearly" element={<LogYearlyRevenue />} />
            <Route path="/dailyOps" element={<DailyOps />} />
            <Route path="/RevenuePie" element={<RevPie />} />
            <Route path="/LogisticsMTD" element={<LogisticsMTD />} />
            <Route path="/EmptyMiles" element={<EmptyMilesWrapper />} />
            <Route path="/Orders" element={<OrdersPage />} />
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
