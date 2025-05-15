import { BrowserRouter, Routes, Route } from 'react-router-dom';
import {Contacts, ContactShow} from './contacts';

const App = () =>{
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Contacts />} />
        <Route path="/contacts" element={<Contacts />} />
        <Route path="/contacts/:id" element={<ContactShow />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
