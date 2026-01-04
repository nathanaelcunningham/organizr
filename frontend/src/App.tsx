
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/layout/Layout';
import { Notifications } from './components/common/Notifications';
import { SearchPage } from './pages/SearchPage';
import { DownloadsPage } from './pages/DownloadsPage';
import {
  ConfigPage,
  ProvidersConfigPage,
  GeneralConfigPage,
} from './pages/ConfigPage';
import { NotFoundPage } from './pages/NotFoundPage';

function App() {
  return (
    <BrowserRouter>
      <Notifications />
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/search" replace />} />
          <Route path="search" element={<SearchPage />} />
          <Route path="downloads" element={<DownloadsPage />} />
          <Route path="config" element={<ConfigPage />}>
            <Route index element={<Navigate to="providers" replace />} />
            <Route path="providers" element={<ProvidersConfigPage />} />
            <Route path="general" element={<GeneralConfigPage />} />
          </Route>
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
