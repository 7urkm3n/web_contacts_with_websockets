import { Fragment, useEffect, useState, useRef, useCallback} from "react";
import { useNavigate } from 'react-router-dom';
import config from './../config'
import AddContact from './add'

export const Main = () => {
  const navigate = useNavigate();
  const [contacts, setContacts] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const ws = useRef(null);

  const connect = useCallback(() => {
    if (ws.current) {
      ws.current.close();
    }

    ws.current = new WebSocket(config.WS_URL);

    ws.current.onopen = () => {
      console.log("ws opened");
      setIsConnected(true);
    };

    ws.current.onclose = () => {
      console.log("ws closed");
      setIsConnected(false);
      // Handle reconnection if needed
    };

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log("TRIGGERED: ", message);
      
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
        setContacts(prev => prev.filter(c => c.id !== message.contact.id));
      }
    };

    ws.current.onerror = (error) => {
       setIsConnected(false);
    };
  }, [config.WS_URL]);


  useEffect(() => {
    connect();

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [connect]);

  return (
    <div className="p-12 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Contacts</h1>
      <div className="grid gap-4">
        <AddContact />
      </div>
      <div className="grid gap-4">
        <div className="relative overflow-x-auto">
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
