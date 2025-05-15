import React, { useState, useEffect, Fragment} from 'react';
import { useLocation } from 'react-router-dom';
import config from './../config'

export const Show = () => {
  const location = useLocation();
  const [history, setHistory] = useState([]);
  const contact = location.state;

  const fetchContactHistory = async () => {
    const res = await fetch(`${config.API_URL}/contacts/${location.state.id}/history`);
    const data = await res.json();
    let changed = [];
    if(data){
      data.map((history) => (
        changed.push(parseJSON(history.changes))
      ))
    }
    setHistory(changed);
  };

  useEffect(() => {
    fetchContactHistory();
  }, [])

  const parseJSON = (changes) => {
    const val = JSON.parse(changes)
    const firstKey = Object.keys(val)[0];
    return val[firstKey]
  }

  return ( 
    <div className="p-6 max-w-4xl mx-auto">
      <div className="bg-white max-w-2xl shadow overflow-hidden sm:rounded-lg">
          <div className="px-4 py-5 sm:px-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900">
                {contact.first_name} {contact.last_name}
              </h3>
              <p className="mt-1 max-w-2xl text-sm text-gray-500">
                <b> Email: </b> {contact.email} 
              </p>
              <p className="mt-1 max-w-2xl text-sm text-gray-500">
                <b> Phone Number: </b> {contact.phone_number} 
              </p>
          </div>
          <div className="border-t border-gray-200">
            <dl>
              <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <h1> <u> Changes history </u> </h1>
                <br/>
                    {history && history.map((h, i) => (     
                      <dd key={i} className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2"> 
                        <b>From: </b> {h?.from} | <b>To: </b> {h?.to}
                      </dd>
                    ))}
                  </div>
              </dl>
          </div>
      </div>
    </div>
  );
}
