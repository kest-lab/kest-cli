import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { User, KeyRound, Shield, Trash2, Loader2, Mail } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/components/ui/alert-dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import { useAuthStore } from '@/store/auth-store'
import { authApi } from '@/services/auth'
import request from '@/http'

// ========== Profile Form ==========

const profileSchema = z.object({
  nickname: z.string().max(50).optional(),
  phone: z.string().max(20).optional(),
  bio: z.string().max(500).optional(),
})

type ProfileFormData = z.infer<typeof profileSchema>

// ========== Password Form ==========

const passwordSchema = z.object({
  old_password: z.string().min(1, 'Current password is required'),
  new_password: z.string().min(6, 'New password must be at least 6 characters'),
  confirm_password: z.string().min(1, 'Please confirm your new password'),
}).refine(data => data.new_password === data.confirm_password, {
  message: "Passwords don't match",
  path: ['confirm_password'],
})

type PasswordFormData = z.infer<typeof passwordSchema>

// ========== Component ==========

export function SettingsPage() {
  const { user, setUser } = useAuthStore()
  const [profileLoading, setProfileLoading] = useState(false)
  const [passwordLoading, setPasswordLoading] = useState(false)
  const [resetLoading, setResetLoading] = useState(false)
  const [resetEmail, setResetEmail] = useState(user?.email || '')
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [deleteConfirmName, setDeleteConfirmName] = useState('')
  const [deleteLoading, setDeleteLoading] = useState(false)

  const profileForm = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      nickname: user?.nickname || '',
      phone: user?.phone || '',
      bio: user?.bio || '',
    },
  })

  const passwordForm = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema),
    defaultValues: { old_password: '', new_password: '', confirm_password: '' },
  })

  const onProfileSubmit = async (data: ProfileFormData) => {
    setProfileLoading(true)
    try {
      await request.put('/v1/users/profile', data)
      const refreshed = await authApi.getProfile()
      setUser(refreshed)
      toast.success('Profile updated')
    } catch (err: any) {
      toast.error(err.message || 'Failed to update profile')
    } finally {
      setProfileLoading(false)
    }
  }

  const onPasswordSubmit = async (data: PasswordFormData) => {
    setPasswordLoading(true)
    try {
      await authApi.changePassword({
        old_password: data.old_password,
        new_password: data.new_password,
      })
      toast.success('Password changed successfully')
      passwordForm.reset()
    } catch (err: any) {
      toast.error(err.message || 'Failed to change password')
    } finally {
      setPasswordLoading(false)
    }
  }

  const onResetPassword = async () => {
    if (!resetEmail) {
      toast.error('Please enter an email address')
      return
    }
    setResetLoading(true)
    try {
      await authApi.resetPassword(resetEmail)
      toast.success('Password reset email sent')
    } catch (err: any) {
      toast.error(err.message || 'Failed to send password reset email')
    } finally {
      setResetLoading(false)
    }
  }

  const onDeleteAccount = async () => {
    setDeleteLoading(true)
    try {
      await request.delete('/v1/users/account')
      toast.success('Account deleted successfully')
      setDeleteDialogOpen(false)
      useAuthStore.getState().clearAuth()
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete account')
    } finally {
      setDeleteLoading(false)
    }
  }

  if (!user) {
    return (
      <div className="container mx-auto p-8 text-center">
        <p className="text-muted-foreground">Please log in to view settings.</p>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-8 max-w-3xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">Manage your account and preferences</p>
      </div>

      <Tabs defaultValue="profile" className="space-y-6">
        <TabsList>
          <TabsTrigger value="profile" className="gap-2">
            <User className="h-4 w-4" /> Profile
          </TabsTrigger>
          <TabsTrigger value="security" className="gap-2">
            <KeyRound className="h-4 w-4" /> Security
          </TabsTrigger>
        </TabsList>

        {/* Profile Tab */}
        <TabsContent value="profile">
          <Card>
            <CardHeader>
              <CardTitle>Profile Information</CardTitle>
              <CardDescription>Update your personal details</CardDescription>
            </CardHeader>
            <CardContent>
              {/* Read-only fields */}
              <div className="space-y-4 mb-6">
                <div className="flex items-center gap-4 p-4 bg-muted/50 rounded-lg">
                  <div className="w-14 h-14 rounded-full bg-gradient-to-br from-blue-600 to-purple-600 flex items-center justify-center text-white text-xl font-bold shrink-0">
                    {user.username?.charAt(0).toUpperCase() || 'U'}
                  </div>
                  <div>
                    <p className="font-semibold text-lg">{user.username}</p>
                    <p className="text-sm text-muted-foreground">{user.email}</p>
                    <Badge variant="outline" className="mt-1 text-xs">
                      ID: {user.id}
                    </Badge>
                  </div>
                </div>
              </div>

              {/* Editable fields */}
              <form onSubmit={profileForm.handleSubmit(onProfileSubmit)} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="nickname">Nickname</Label>
                    <Input id="nickname" placeholder="Display name" {...profileForm.register('nickname')} />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="phone">Phone</Label>
                    <Input id="phone" placeholder="+1 234 567 890" {...profileForm.register('phone')} />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="bio">Bio</Label>
                  <Textarea
                    id="bio"
                    placeholder="Tell us about yourself..."
                    rows={3}
                    className="resize-none"
                    {...profileForm.register('bio')}
                  />
                </div>
                <div className="flex justify-end pt-2">
                  <Button type="submit" disabled={!profileForm.formState.isDirty || profileLoading}>
                    {profileLoading ? <><Loader2 className="h-4 w-4 mr-2 animate-spin" /> Saving...</> : 'Save Changes'}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Tab */}
        <TabsContent value="security" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Change Password</CardTitle>
              <CardDescription>Update your password to keep your account secure</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={passwordForm.handleSubmit(onPasswordSubmit)} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="old_password">Current Password</Label>
                  <Input id="old_password" type="password" {...passwordForm.register('old_password')} />
                  {passwordForm.formState.errors.old_password && (
                    <p className="text-sm text-red-500">{passwordForm.formState.errors.old_password.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="new_password">New Password</Label>
                  <Input id="new_password" type="password" {...passwordForm.register('new_password')} />
                  {passwordForm.formState.errors.new_password && (
                    <p className="text-sm text-red-500">{passwordForm.formState.errors.new_password.message}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirm_password">Confirm New Password</Label>
                  <Input id="confirm_password" type="password" {...passwordForm.register('confirm_password')} />
                  {passwordForm.formState.errors.confirm_password && (
                    <p className="text-sm text-red-500">{passwordForm.formState.errors.confirm_password.message}</p>
                  )}
                </div>
                <div className="flex justify-end pt-2">
                  <Button type="submit" disabled={passwordLoading}>
                    {passwordLoading ? <><Loader2 className="h-4 w-4 mr-2 animate-spin" /> Changing...</> : 'Change Password'}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Mail className="h-5 w-5" /> Reset Password
              </CardTitle>
              <CardDescription>Send a password reset link to an email address</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex gap-3 items-end">
                <div className="flex-1 space-y-2">
                  <Label htmlFor="reset_email">Email Address</Label>
                  <Input
                    id="reset_email"
                    type="email"
                    placeholder="user@example.com"
                    value={resetEmail}
                    onChange={e => setResetEmail(e.target.value)}
                  />
                </div>
                <Button onClick={onResetPassword} disabled={resetLoading || !resetEmail}>
                  {resetLoading ? <><Loader2 className="h-4 w-4 mr-2 animate-spin" /> Sending...</> : 'Send Reset Email'}
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card className="border-red-200">
            <CardHeader>
              <CardTitle className="text-red-600 flex items-center gap-2">
                <Shield className="h-5 w-5" /> Danger Zone
              </CardTitle>
              <CardDescription>Irreversible actions for your account</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between p-4 border border-red-200 rounded-lg bg-red-50 dark:bg-red-950/20">
                <div>
                  <h4 className="font-semibold text-red-900 dark:text-red-400">Delete Account</h4>
                  <p className="text-sm text-red-700 dark:text-red-500">Permanently delete your account and all associated data</p>
                </div>
                <AlertDialog open={deleteDialogOpen} onOpenChange={(open) => { setDeleteDialogOpen(open); if (!open) setDeleteConfirmName('') }}>
                  <AlertDialogTrigger asChild>
                    <Button variant="destructive" size="sm">
                      <Trash2 className="h-4 w-4 mr-2" /> Delete
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Delete Account</AlertDialogTitle>
                      <AlertDialogDescription>
                        This action is permanent and cannot be undone. All your data will be deleted.
                        Please type <span className="font-semibold text-foreground">{user.username}</span> to confirm.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <div className="py-2">
                      <Input
                        placeholder="Enter your username"
                        value={deleteConfirmName}
                        onChange={e => setDeleteConfirmName(e.target.value)}
                      />
                    </div>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <Button
                        variant="destructive"
                        disabled={deleteConfirmName !== user.username || deleteLoading}
                        onClick={onDeleteAccount}
                      >
                        {deleteLoading ? <><Loader2 className="h-4 w-4 mr-2 animate-spin" /> Deleting...</> : 'Delete Account'}
                      </Button>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
