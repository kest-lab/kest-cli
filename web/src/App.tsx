import { Routes, Route, Navigate } from 'react-router-dom';
import { ProjectsPage } from '@/pages/projects';
import { ProjectDetailPage } from '@/pages/projects/detail';
import { FlowEditorPage } from '@/pages/projects/flow-editor';
import { HomePage } from '@/pages/home';
import { LoginPage } from '@/pages/auth/login';
import { RegisterPage } from '@/pages/auth/register';
import { SettingsPage } from '@/pages/settings';
import { AdminLayout } from '@/components/layout/admin-layout';

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/projects" element={<AdminLayout><ProjectsPage /></AdminLayout>} />
      <Route path="/projects/:id" element={<AdminLayout><ProjectDetailPage /></AdminLayout>} />
      <Route path="/projects/:id/api-specs/:sid" element={<AdminLayout><ProjectDetailPage /></AdminLayout>} />
      <Route path="/projects/:id/flows/:fid" element={<AdminLayout><FlowEditorPage /></AdminLayout>} />
      <Route path="/settings" element={<AdminLayout><SettingsPage /></AdminLayout>} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default App;
