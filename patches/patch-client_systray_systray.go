$NetBSD: patch-client_systray_systray.go,v 1.1 2026/07/29 00:00:00 schmonz Exp $

Menu-bar icon only (no "tailscale" text) on macOS, and make the systray
"Connect" action start interactive login on a logged-out node so the
login URL opens in the system browser (as "tailscale up" would), instead
of silently doing nothing.

--- client/systray/systray.go.orig	2026-07-02 18:59:26.000000000 +0000
+++ client/systray/systray.go
@@ -174,7 +174,11 @@
 
 	// set initial title, which is used by the systray package as the ID of the StatusNotifierItem.
 	// This value will get overwritten later as the client status changes.
-	systray.SetTitle("tailscale")
+	// On macOS there is no StatusNotifierItem and the title renders as literal
+	// text next to the menu-bar icon, so skip it there and show the icon only.
+	if runtime.GOOS != "darwin" {
+		systray.SetTitle("tailscale")
+	}
 
 	menu.rebuild()
 
@@ -438,6 +442,17 @@
 			if err != nil {
 				log.Printf("error connecting: %v", err)
 			}
+			// Setting WantRunning alone won't authenticate a logged-out node, so
+			// (like "tailscale up") kick off interactive login when the backend
+			// needs it. That makes tailscaled emit a BrowseToURL, which
+			// watchIPNBus opens in the system browser. Without this, Connect
+			// silently does nothing on a fresh node until the user runs
+			// "tailscale up" by hand.
+			if menu.status != nil && (menu.status.BackendState == "NeedsLogin" || menu.status.BackendState == "NoState") {
+				if err := menu.lc.StartLoginInteractive(ctx); err != nil {
+					log.Printf("error starting interactive login: %v", err)
+				}
+			}
 
 		case <-menu.disconnect.ClickedCh:
 			_, err := menu.lc.EditPrefs(ctx, &ipn.MaskedPrefs{
