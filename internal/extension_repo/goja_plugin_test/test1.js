function init() {
    $ui.register((ctx) => {
        const tray = ctx.newTray({
            tooltipText: "Test Plugin",
            iconUrl: "https://raw.githubusercontent.com/5rahim/hibike/main/icons/seadex.png",
        });

        const currentMediaId = ctx.state(0);
        const storageBackgroundImage = ctx.state("");
        const mediaIds = ctx.state([]);

        const customBannerImageRef = ctx.registerFieldRef("customBannerImageRef");

        const fetchBackgroundImage = () => {
            const backgroundImage = $storage.get('backgroundImages.' + currentMediaId.get());
            if (backgroundImage) {
                storageBackgroundImage.set(backgroundImage);
                customBannerImageRef.setValue(backgroundImage);
            } else {
                storageBackgroundImage.set("");
                customBannerImageRef.setValue("");
            }
        }

        ctx.effect(() => {
            console.log("media ID changed, fetching background image and updating tray");
            fetchBackgroundImage();

            console.log("updating tray");
        }, [currentMediaId]);


        fetchBackgroundImage()

        ctx.screen.onNavigate((e) => {
            console.log("screen navigated", e);
            if (e.pathname === "/entry" && !!e.query) {
                const id = parseInt(e.query.replace("?id=", ""));
                currentMediaId.set(id);
            } else {
                currentMediaId.set(0);
            }

            console.log("updating tray");
        });

        ctx.registerEventHandler("saveBackgroundImage", () => {
            ctx.toast.info("Setting background image to " + customBannerImageRef.current);
            $storage.set('backgroundImages.' + currentMediaId.get(), customBannerImageRef.current);
            ctx.toast.success("Background image saved");
            fetchBackgroundImage();
            $anilist.refreshAnimeCollection();
        });

        // $store.watch("mediaIds", (mId) => {
        // 	mediaIds.set(p => [...p, mId]);
        // });

        ctx.registerEventHandler("button-clicked", () => {
            console.log("button-clicked");
            console.log("navigating to /entry?id=21");
            try {
                ctx.screen.navigateTo("/entry?id=21");
            } catch (e) {
                console.error("navigate error", e);
            }
            ctx.setTimeout(() => {
                try {
                    console.log("navigating to /entry?id=177709");
                    ctx.screen.navigateTo("/entry?id=177709");
                } catch (e) {
                    console.error("navigate error", e);
                }
            }, 1000);
            ctx.setTimeout(() => {
                try {
                    console.log("opening https://google.com");
                    const cmd = $os.cmd("open", "https://google.com");
                    cmd.run();
                } catch (e) {
                    console.error("open error", e);
                }
            }, 2000);
        });

        tray.render(() => {
            return tray.stack({
                items: [
                    tray.button("Click me", {onClick: "button-clicked"}),
                    currentMediaId.get() === 0 ? tray.text("Open an anime or manga") : tray.stack({
                        items: [
                            tray.text(`Current media ID: ${currentMediaId.get()}`),
                            tray.input({fieldRef: "customBannerImageRef", value: storageBackgroundImage.get()}),
                            tray.button({label: "Save", onClick: "saveBackgroundImage"}),
                        ],
                    }),
                ],
            });
        });
    })

    $app.onGetAnime((e) => {
        $store.set("mediaIds", e.anime.id);
        e.next();
    });


    $app.onGetAnimeCollection((e) => {
        const bannerImages = $storage.get('backgroundImages');
        for (let i = 0; i < e.animeCollection.mediaListCollection.lists.length; i++) {
            for (let j = 0; j < e.animeCollection.mediaListCollection.lists[i].entries.length; j++) {
                const mediaId = e.animeCollection.mediaListCollection.lists[i].entries[j].media.id;
                const bannerImage = bannerImages[mediaId.toString()] || "";
                if (!!bannerImage) {
                    $replace(e.animeCollection.mediaListCollection.lists[i].entries[j].media.bannerImage, bannerImage);
                }
            }
        }
        e.next();
    });

    $app.onGetRawAnimeCollection((e) => {
        const bannerImages = $storage.get('backgroundImages');
        for (let i = 0; i < e.animeCollection.mediaListCollection.lists.length; i++) {
            for (let j = 0; j < e.animeCollection.mediaListCollection.lists[i].entries.length; j++) {
                const mediaId = e.animeCollection.mediaListCollection.lists[i].entries[j].media.id;
                const bannerImage = bannerImages[mediaId.toString()] || "";
                if (!!bannerImage) {
                    $replace(e.animeCollection.mediaListCollection.lists[i].entries[j].media.bannerImage, bannerImage);
                }
            }
        }
        e.next();
    });

}
