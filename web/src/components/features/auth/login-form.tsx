'use client';

import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import * as z from 'zod';
import { Eye, EyeOff, Loader2 } from 'lucide-react';

import { cn } from '@/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Google, Apple, GitHub } from '@/components/ui/icons';
import { useLogin } from '@/hooks/use-auth';
import { mockUsers } from '@/config/auth';
import { useAuthStore, authSelectors } from '@/store/auth-store';
import { useT } from '@/i18n';

const loginSchema = z.object({
    email: z.string().min(1, 'emailRequired').email('emailInvalid'),
    password: z.string().min(6, 'passwordTooShort'),
});

type LoginFormData = z.infer<typeof loginSchema>;

export function LoginForm({ className }: { className?: string }) {
    const t = useT();
    const [showPassword, setShowPassword] = useState(false);
    const { mutateAsync: login, isPending: isMutationLoading } = useLogin();
    const systemFeatures = useAuthStore.use.systemFeatures();

    useEffect(() => {
        if (import.meta.env.DEV) {
            console.log('📧 Mock accounts available for login:');
            mockUsers.forEach((user) => {
                console.log(`   Email: ${user.email}, Password: ${user.password}`);
            });
        }
    }, []);

    const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
    });

    const onSubmit = async (data: LoginFormData) => {
        try { await login(data); } catch (err) { console.error(err); }
    };

    const isFormLoading = isMutationLoading || isSubmitting;
    const hasSocial = systemFeatures ? authSelectors.hasSocialLogin(useAuthStore.getState()) : false;

    const getErrorMessage = (errorKey: string | undefined) => {
        if (!errorKey) return undefined;
        return t(`auth.${errorKey}`) || errorKey;
    };

    return (
        <div className={cn('flex flex-col gap-4', className)}>
            <Card className="border-border/60 shadow-xl dark:shadow-2xl dark:shadow-black/20 bg-card/95 backdrop-blur-sm">
                <CardHeader className="text-center space-y-1 px-6 pb-1 pt-5">
                    <CardTitle className="text-2xl font-bold">
                        {t('auth.welcomeBack')}
                    </CardTitle>
                    <p className="text-sm text-muted-foreground">{t('auth.signInToContinue')}</p>
                </CardHeader>

                <CardContent className="space-y-4 px-6 pb-5 pt-3">
                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
                        <div className="space-y-1">
                            <Label htmlFor="email" className="text-sm font-medium">{t('auth.email')}</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder={t('auth.enterEmail')}
                                disabled={isFormLoading}
                                className="h-10 bg-muted/30 border-border/80 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
                                {...register('email')}
                            />
                            {errors.email && (
                                <p className="text-xs text-destructive">{getErrorMessage(errors.email.message)}</p>
                            )}
                        </div>

                        <div className="space-y-1">
                            <div className="flex items-center justify-between">
                                <Label htmlFor="password" className="text-sm font-medium">{t('auth.password')}</Label>
                                <Link to="/forgot-password" title={t('auth.forgotPassword')} className="text-primary text-xs font-medium hover:underline transition-colors">
                                    {t('auth.forgotPassword')}
                                </Link>
                            </div>
                            <div className="relative">
                                <Input
                                    id="password"
                                    type={showPassword ? 'text' : 'password'}
                                    placeholder={t('auth.enterPassword')}
                                    disabled={isFormLoading}
                                    className="h-10 pr-10 bg-muted/30 border-border/80 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all"
                                    {...register('password')}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={() => setShowPassword(!showPassword)}
                                    tabIndex={-1}
                                >
                                    {showPassword ? (
                                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                        <Eye className="h-4 w-4 text-muted-foreground" />
                                    )}
                                </Button>
                            </div>
                            {errors.password && (
                                <p className="text-xs text-destructive">{getErrorMessage(errors.password.message)}</p>
                            )}
                        </div>

                        <Button
                            type="submit"
                            className="w-full h-10 font-medium bg-linear-to-r from-primary to-primary-deeper hover:from-primary/90 hover:to-primary-deeper/90 shadow-md hover:shadow-lg transition-all duration-200"
                            disabled={isFormLoading}
                        >
                            {isFormLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('auth.signIn')}
                        </Button>
                    </form>

                    {hasSocial && (
                        <>
                            <div className="relative">
                                <div className="absolute inset-0 flex items-center">
                                    <span className="w-full border-t border-border/50" />
                                </div>
                                <div className="relative flex justify-center text-xs">
                                    <span className="bg-card px-3 text-muted-foreground">
                                        {t('auth.orContinueWith')}
                                    </span>
                                </div>
                            </div>

                            <div className="grid grid-cols-3 gap-2">
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}>
                                    <Google className="h-4 w-4" />
                                </Button>
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}>
                                    <Apple className="h-4 w-4" />
                                </Button>
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}>
                                    <GitHub className="h-4 w-4" />
                                </Button>
                            </div>
                        </>
                    )}

                    <div className="text-center text-sm text-muted-foreground">
                        {t('auth.noAccount')}{' '}
                        <Link to="/register" title={t('auth.signUp')} className="font-semibold text-primary hover:underline transition-colors">
                            {t('auth.signUp')}
                        </Link>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
