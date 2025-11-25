'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { GoogleLogin, GoogleOAuthProvider } from '@react-oauth/google';
import { useAuthStore } from '@/shared/lib';
import { AuthAPI } from '@/shared/lib/api/auth';
import { Mail, Lock, Eye, EyeOff, ArrowLeft } from 'lucide-react';
import Link from 'next/link';

const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || '';

export default function LoginPage() {
  const router = useRouter();
  const { isAuthenticated, setAuth, setLoading } = useAuthStore();
  const [isSignup, setIsSignup] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [error, setError] = useState('');
  const [returnUrl, setReturnUrl] = useState<string | null>(null);

  // Password reset state
  const [showResetPassword, setShowResetPassword] = useState(false);
  const [resetEmail, setResetEmail] = useState('');
  const [resetToken, setResetToken] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [resetStep, setResetStep] = useState<'request' | 'reset'>('request');
  const [resetSuccess, setResetSuccess] = useState('');

  // Get return URL from query params on mount
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const from = params.get('from');
    if (from) {
      setReturnUrl(from);
    }
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      // Redirect to return URL if specified, otherwise to chat
      router.push(returnUrl || '/chat');
    }
  }, [isAuthenticated, router, returnUrl]);

  const handleGoogleSuccess = async (credentialResponse: any) => {
    if (!credentialResponse.credential) {
      console.error('No credential in response');
      return;
    }

    setLoading(true);

    try {
      const authResponse = await AuthAPI.googleLogin(credentialResponse.credential);

      setAuth(authResponse.user, {
        access_token: authResponse.access_token,
        refresh_token: authResponse.refresh_token,
        expires_in: authResponse.expires_in,
      });

      // Note: We don't sync preferences or load messages here
      // The use-chat.ts hook will handle this automatically on mount for authenticated users
      // This prevents race conditions with localStorage persistence
      router.push(returnUrl || '/chat');
    } catch (error) {
      // Log detailed error in development only
      if (process.env.NODE_ENV === 'development') {
        console.error('Login failed:', error);
      }
      setError(error instanceof Error ? error.message : 'Failed to login with Google');
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleError = () => {
    // Log detailed error in development only
    if (process.env.NODE_ENV === 'development') {
      console.error('Google Login Failed');
    }
    setError('Failed to login with Google');
  };

  const handleEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const authResponse = isSignup
        ? await AuthAPI.signup(email, password, fullName)
        : await AuthAPI.login(email, password);

      setAuth(authResponse.user, {
        access_token: authResponse.access_token,
        refresh_token: authResponse.refresh_token,
        expires_in: authResponse.expires_in,
      });

      // Note: We don't sync preferences or load messages here
      // The use-chat.ts hook will handle this automatically on mount for authenticated users
      // This prevents race conditions with localStorage persistence
      router.push(returnUrl || '/chat');
    } catch (error) {
      // Log detailed error in development only
      if (process.env.NODE_ENV === 'development') {
        console.error('Auth failed:', error);
      }
      setError(error instanceof Error ? error.message : 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  const handleRequestPasswordReset = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setResetSuccess('');
    setLoading(true);

    try {
      const response = await AuthAPI.requestPasswordReset(resetEmail);
      // In production, the token would be sent via email
      // For testing, we receive it in the response
      if (response.token) {
        setResetToken(response.token);
        setResetStep('reset');
        setResetSuccess('Reset token generated. In production, this would be sent to your email.');
      } else {
        setResetSuccess('If an account exists with this email, a password reset link has been sent.');
      }
    } catch (error) {
      // Log detailed error in development only
      if (process.env.NODE_ENV === 'development') {
        console.error('Password reset request failed:', error);
      }
      setError(error instanceof Error ? error.message : 'Failed to request password reset');
    } finally {
      setLoading(false);
    }
  };

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setResetSuccess('');
    setLoading(true);

    try {
      await AuthAPI.resetPassword(resetToken, newPassword);
      setResetSuccess('Password reset successfully! You can now login with your new password.');
      setTimeout(() => {
        setShowResetPassword(false);
        setResetStep('request');
        setResetEmail('');
        setResetToken('');
        setNewPassword('');
        setError('');
        setResetSuccess('');
      }, 2000);
    } catch (error) {
      // Log detailed error in development only
      if (process.env.NODE_ENV === 'development') {
        console.error('Password reset failed:', error);
      }
      setError(error instanceof Error ? error.message : 'Failed to reset password');
    } finally {
      setLoading(false);
    }
  };

  if (isAuthenticated) {
    return null;
  }

  return (
    <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
      <div className="min-h-screen flex flex-col bg-background">
        {/* Header with back button */}
        <header className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-sm mobile-opaque-bg border-b border-border">
          <div className="container mx-auto px-4 h-16 flex items-center">
            <Link
              href="/"
              className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            >
              <ArrowLeft className="w-5 h-5" />
              <span className="text-sm font-medium">Back</span>
            </Link>
          </div>
        </header>

        {/* Main content */}
        <div className="flex-1 flex items-center justify-center px-4 pt-24 pb-8">
          <div className="w-full max-w-md space-y-8">
            {/* Password Reset Flow */}
            {showResetPassword ? (
              <>
                {/* Reset Password Title */}
                <div className="text-center space-y-3">
                  <h1 className="text-4xl font-semibold tracking-tight">
                    {resetStep === 'request' ? 'Reset Password' : 'Set New Password'}
                  </h1>
                  <p className="text-muted-foreground">
                    {resetStep === 'request'
                      ? 'Enter your email to receive a reset link'
                      : 'Enter your new password'}
                  </p>
                </div>

                {/* Success message */}
                {resetSuccess && (
                  <div className="p-4 rounded-lg bg-green-50 dark:bg-green-900/10 border border-green-200 dark:border-green-800/30">
                    <p className="text-sm text-green-800 dark:text-green-200">{resetSuccess}</p>
                  </div>
                )}

                {/* Error message */}
                {error && (
                  <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/10 border border-red-200 dark:border-red-800/30">
                    <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
                  </div>
                )}

                {resetStep === 'request' ? (
                  /* Request Reset Token Form */
                  <form onSubmit={handleRequestPasswordReset} className="space-y-4">
                    <div className="space-y-2">
                      <label htmlFor="reset-email" className="block text-sm font-medium">
                        Email
                      </label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <input
                          id="reset-email"
                          type="email"
                          value={resetEmail}
                          onChange={(e) => setResetEmail(e.target.value)}
                          required
                          className="w-full pl-10 pr-3 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                          placeholder="you@example.com"
                        />
                      </div>
                    </div>

                    <button
                      type="submit"
                      className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium py-2.5 rounded-lg transition-colors cursor-pointer"
                    >
                      Send Reset Link
                    </button>

                    <div className="text-center">
                      <button
                        type="button"
                        onClick={() => {
                          setShowResetPassword(false);
                          setError('');
                          setResetSuccess('');
                          setResetEmail('');
                        }}
                        className="text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                      >
                        Back to login
                      </button>
                    </div>
                  </form>
                ) : (
                  /* Reset Password Form */
                  <form onSubmit={handleResetPassword} className="space-y-4">
                    <div className="space-y-2">
                      <label htmlFor="reset-token" className="block text-sm font-medium">
                        Reset Token
                      </label>
                      <input
                        id="reset-token"
                        type="text"
                        value={resetToken}
                        onChange={(e) => setResetToken(e.target.value)}
                        required
                        className="w-full px-3 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                        placeholder="Enter reset token"
                      />
                      <p className="text-xs text-muted-foreground">
                        In production, this would be sent to your email
                      </p>
                    </div>

                    <div className="space-y-2">
                      <label htmlFor="new-password-reset" className="block text-sm font-medium">
                        New Password
                      </label>
                      <div className="relative">
                        <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <input
                          id="new-password-reset"
                          type={showNewPassword ? 'text' : 'password'}
                          value={newPassword}
                          onChange={(e) => setNewPassword(e.target.value)}
                          required
                          minLength={8}
                          className="w-full pl-10 pr-10 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                          placeholder="••••••••"
                        />
                        <button
                          type="button"
                          onClick={() => setShowNewPassword(!showNewPassword)}
                          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                        >
                          {showNewPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                        </button>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        Must be at least 8 characters
                      </p>
                    </div>

                    <button
                      type="submit"
                      className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium py-2.5 rounded-lg transition-colors cursor-pointer"
                    >
                      Reset Password
                    </button>

                    <div className="text-center">
                      <button
                        type="button"
                        onClick={() => {
                          setResetStep('request');
                          setResetToken('');
                          setNewPassword('');
                          setError('');
                        }}
                        className="text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                      >
                        Back
                      </button>
                    </div>
                  </form>
                )}
              </>
            ) : (
              /* Login/Signup Flow */
              <>
                {/* Logo and title */}
                <div className="text-center space-y-3">
                  <h1 className="text-4xl font-semibold tracking-tight">
                    {isSignup ? 'Create your account' : 'Welcome back'}
                  </h1>
                  <p className="text-muted-foreground">
                    {isSignup
                      ? 'Start finding the best deals today'
                      : 'Continue to your account'}
                  </p>
                </div>

            {/* Error message */}
            {error && (
              <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/10 border border-red-200 dark:border-red-800/30">
                <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
              </div>
            )}

            {/* Google login */}
            <div className="space-y-4">
              <div className="flex justify-center">
                <GoogleLogin
                  onSuccess={handleGoogleSuccess}
                  onError={handleGoogleError}
                  theme="outline"
                  size="large"
                  text={isSignup ? "signup_with" : "signin_with"}
                  shape="rectangular"
                  width="384"
                />
              </div>

              {/* Divider */}
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-border"></div>
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-background px-2 text-muted-foreground">
                    Or continue with email
                  </span>
                </div>
              </div>
            </div>

            {/* Email/Password form */}
            <form onSubmit={handleEmailSubmit} className="space-y-4">
              {isSignup && (
                <div className="space-y-2">
                  <label
                    htmlFor="fullName"
                    className="block text-sm font-medium"
                  >
                    Full Name
                  </label>
                  <input
                    id="fullName"
                    type="text"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                    className="w-full px-3 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                    placeholder="Enter your name"
                    required={isSignup}
                  />
                </div>
              )}

              <div className="space-y-2">
                <label
                  htmlFor="email"
                  className="block text-sm font-medium"
                >
                  Email
                </label>
                <div className="relative">
                  <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full pl-10 pr-3 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                    placeholder="you@example.com"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <label
                    htmlFor="password"
                    className="block text-sm font-medium"
                  >
                    Password
                  </label>
                  {!isSignup && (
                    <button
                      type="button"
                      onClick={() => {
                        setShowResetPassword(true);
                        setError('');
                      }}
                      className="text-xs text-primary hover:text-primary/80 transition-colors cursor-pointer"
                    >
                      Forgot password?
                    </button>
                  )}
                </div>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    minLength={8}
                    className="w-full pl-10 pr-10 py-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
                {isSignup && (
                  <p className="text-xs text-muted-foreground">
                    Must be at least 8 characters
                  </p>
                )}
              </div>

              <button
                type="submit"
                className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium py-2.5 rounded-lg transition-colors cursor-pointer"
              >
                {isSignup ? 'Create account' : 'Continue'}
              </button>
            </form>

            {/* Toggle sign up/sign in */}
            <div className="text-center">
              <button
                onClick={() => {
                  setIsSignup(!isSignup);
                  setError('');
                }}
                className="text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              >
                {isSignup
                  ? 'Already have an account? '
                  : "Don't have an account? "}
                <span className="font-medium text-primary">
                  {isSignup ? 'Sign in' : 'Sign up'}
                </span>
              </button>
            </div>

            {/* Terms */}
            <p className="text-center text-xs text-muted-foreground">
              By continuing, you agree to our{' '}
              <Link href="/terms-of-use" className="underline hover:text-foreground">
                Terms of Service
              </Link>
              {' '}and{' '}
              <Link href="/privacy-policy" className="underline hover:text-foreground">
                Privacy Policy
              </Link>
            </p>
              </>
            )}
          </div>
        </div>
      </div>
    </GoogleOAuthProvider>
  );
}
