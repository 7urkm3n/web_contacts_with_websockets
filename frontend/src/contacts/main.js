import { Fragment, useEffect, useState, useRef } from "react";
import { useNavigate, Link } from 'react-router-dom';
import config from './../config'

export const Main = () => {
  const navigate = useNavigate();
  const [contacts, setContacts] = useState([]);
  const [contact, setContact] = useState();

  // const fetchContacts = async () => {
  //   const res = await fetch(`${config.API_URL}/contacts`);
  //   const data = await res.json();
  //   setContacts(data);
  // };

  // useEffect(() => {
  //   fetchContacts();
  // }, []);

  useEffect(() => {
    const ws = new WebSocket(config.WS_URL);
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'allContacts') {
        setContacts(message.contacts);
      } else if (message.type === 'createContact') {
        setContacts(prevContacts => [...prevContacts, message.contact]);
      } else if (message.type === 'updateContact') {
        setContacts(prevContacts => 
          prevContacts.map( contact => 
            contact.id === message.contact.id ? message.contact : contact 
        ));
      } else if (message.type === 'deleteContact') {
        console.log("TRIGGERED: ", message);
        setContacts(prev => prev.filter(c => c.id !== message.contact.id));
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  const createContact = () => {
    console.log('contacts :>> ', contacts);
    // console.log("createContact: ", contact);
    
    // let obj = {
    //   first_name: "ro",
    //   last_name: "go",
    //   email: "rr@asdb.co",
    //   phone_number: "1234423123"
    // }

    // const ws = new WebSocket('ws://localhost:4000/ws');
    // const c = {
    //   contact: obj,
    // };

    // ws.onopen = () => {
    //   console.log("OBJ: ", obj);
    //   ws.send(JSON.stringify(obj));
    // };
    // setContact('');
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Contacts</h1>
      <div className="grid gap-4">
        <div className="relative overflow-x-auto">
          
          <button onClick={() => createContact()}> fetch </button>

          <table className="w-full text-sm text-left rtl:text-right text-gray-500 ">
              <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                  <tr>
                      <th scope="col" className="px-6 py-3">
                          Firstname
                      </th>
                      <th scope="col" className="px-6 py-3">
                          Lastname
                      </th>
                      <th scope="col" className="px-6 py-3">
                          Email
                      </th>
                      <th scope="col" className="px-6 py-3">
                          Phone Number
                      </th>
                  </tr>
              </thead>
              <tbody>
              {contacts && contacts.map((contact) => (
                <Fragment key={contact.id}>
                  <tr className="bg-white border-b" onClick={() => navigate(`/contacts/${contact.id}`, { state: contact }) }>
                      <td scope="row" className="px-6 py-4">
                          {contact.first_name}
                      </td>
                      <td className="px-6 py-4">
                          {contact.last_name}
                      </td>
                      <td className="px-6 py-4">
                          {contact.email}
                      </td>
                      <td className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap">
                          {contact.phone_number}
                      </td>
                  </tr>
                  </Fragment>
              ))}
              </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
