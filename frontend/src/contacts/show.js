import React, { useState, useEffect, Fragment} from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import config from './../config'

export const Show = () => {
  const location = useLocation();
  const navigate = useNavigate();
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

  const deleteContact = async (id) => {
    console.log("id: ", id);

      fetch(`${config.API_URL}/contacts/${id}`, {
        headers: {
          'Content-Type': 'application/json'
        },
        method: 'DELETE',
        mode: 'cors'
      })
      .then(
        navigate("/contacts")
      )
      .catch(console.log("Cant delete"));
  }

  const editContact = async (id) => {
    console.log("Edit Contact Implement");
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

          <div className="px-4">
            <button onClick={()=> editContact} className="text-white inline-flex items-center bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
              Edit
            </button>
            <span>&nbsp;&nbsp;</span>
            <button onClick={()=> deleteContact(contact.id)} className="text-white inline-flex items-center bg-red-700 hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-800">
              Delete
            </button>
          </div>

          <div className="px-4">
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
