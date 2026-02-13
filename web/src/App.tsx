import { Routes, Route, Navigate } from 'react-router-dom';
import { ProjectsPage } from '@/pages/projects';
import { ProjectDetailPage } from '@/pages/projects/detail';
import { APISpecDetailPage } from '@/pages/projects/api-spec-detail';
import { FlowEditorPage } from '@/pages/projects/flow-editor';
import { HomePage } from '@/pages/home';
import { LoginPage } from '@/pages/auth/login';
import { RegisterPage } from '@/pages/auth/register';
import { SettingsPage } from '@/pages/settings';
import { AdminLayout } from '@/components/layout/admin-layout';

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/*" element={
        <AdminLayout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="/projects/:id" element={<ProjectDetailPage />} />
            <Route path="/projects/:id/api-specs/:sid" element={<ProjectDetailPage />} />
            <Route path="/projects/:id/flows/:fid" element={<FlowEditorPage />} />
            <Route path="/settings" element={<SettingsPage />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </AdminLayout>
      } />
    </Routes>
  );
}

export default App;
