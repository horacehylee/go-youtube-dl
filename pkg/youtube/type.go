package youtube

import "fmt"

// StreamFormat from youtube video info
type StreamFormat struct {
	Itag            int    `json:"itag"`
	URL             string `json:"url"`
	MimeType        string `json:"mimeType"`
	Quality         string `json:"quality"`
	ContentLength   string `json:"contentLength"`
	AudioQuality    string `json:"audioQuality"`
	AverageBitrate  int    `json:"averageBitrate"`
	SignatureCipher string `json:"signatureCipher"`
}

func (s StreamFormat) getURL(videoID string) (string, error) {
	if s.URL == "" && s.SignatureCipher == "" {
		return "", fmt.Errorf("Both url and signature cipher is empty")
	}
	if s.URL != "" {
		return s.URL, nil
	}
	return decryptCipher(videoID, s.SignatureCipher)
}

// PlayerResponse from youtube video info
type PlayerResponse struct {
	ResponseContext struct {
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		WebResponseContextExtensionData struct {
			HasDecorated bool `json:"hasDecorated"`
		} `json:"webResponseContextExtensionData"`
	} `json:"responseContext"`
	PlayabilityStatus struct {
		Status          string `json:"status"`
		PlayableInEmbed bool   `json:"playableInEmbed"`
		Miniplayer      struct {
			MiniplayerRenderer struct {
				PlaybackMode string `json:"playbackMode"`
			} `json:"miniplayerRenderer"`
		} `json:"miniplayer"`
		ContextParams string `json:"contextParams"`
	} `json:"playabilityStatus"`
	StreamingData struct {
		ExpiresInSeconds string         `json:"expiresInSeconds"`
		Formats          []StreamFormat `json:"formats"`
		AdaptiveFormats  []StreamFormat `json:"adaptiveFormats"`
		ProbeURL         string         `json:"probeUrl"`
	} `json:"streamingData"`
	PlaybackTracking struct {
		VideostatsPlaybackURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsPlaybackUrl"`
		VideostatsDelayplayURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsDelayplayUrl"`
		VideostatsWatchtimeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsWatchtimeUrl"`
		PtrackingURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"ptrackingUrl"`
		QoeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"qoeUrl"`
		SetAwesomeURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"setAwesomeUrl"`
		AtrURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"atrUrl"`
		YoutubeRemarketingURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"youtubeRemarketingUrl"`
	} `json:"playbackTracking"`
	Captions struct {
		PlayerCaptionsRenderer struct {
			BaseURL    string `json:"baseUrl"`
			Visibility string `json:"visibility"`
			Contribute struct {
				CaptionsMetadataRenderer struct {
					AddSubtitlesText struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								URLEndpoint struct {
									URL string `json:"url"`
								} `json:"urlEndpoint"`
							} `json:"navigationEndpoint"`
						} `json:"runs"`
					} `json:"addSubtitlesText"`
					NoSubtitlesText struct {
						SimpleText string `json:"simpleText"`
					} `json:"noSubtitlesText"`
					PromoSubtitlesText struct {
						SimpleText string `json:"simpleText"`
					} `json:"promoSubtitlesText"`
				} `json:"captionsMetadataRenderer"`
			} `json:"contribute"`
		} `json:"playerCaptionsRenderer"`
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []struct {
				BaseURL string `json:"baseUrl"`
				Name    struct {
					SimpleText string `json:"simpleText"`
				} `json:"name"`
				VssID          string `json:"vssId"`
				LanguageCode   string `json:"languageCode"`
				IsTranslatable bool   `json:"isTranslatable"`
			} `json:"captionTracks"`
			AudioTracks []struct {
				CaptionTrackIndices []int `json:"captionTrackIndices"`
			} `json:"audioTracks"`
			TranslationLanguages []struct {
				LanguageCode string `json:"languageCode"`
				LanguageName struct {
					SimpleText string `json:"simpleText"`
				} `json:"languageName"`
			} `json:"translationLanguages"`
			DefaultAudioTrackIndex int `json:"defaultAudioTrackIndex"`
			Contribute             struct {
				CaptionsMetadataRenderer struct {
					AddSubtitlesText struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								URLEndpoint struct {
									URL string `json:"url"`
								} `json:"urlEndpoint"`
							} `json:"navigationEndpoint"`
						} `json:"runs"`
					} `json:"addSubtitlesText"`
					NoSubtitlesText struct {
						SimpleText string `json:"simpleText"`
					} `json:"noSubtitlesText"`
					PromoSubtitlesText struct {
						SimpleText string `json:"simpleText"`
					} `json:"promoSubtitlesText"`
				} `json:"captionsMetadataRenderer"`
			} `json:"contribute"`
		} `json:"playerCaptionsTracklistRenderer"`
	} `json:"captions"`
	VideoDetails struct {
		VideoID          string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		IsOwnerViewing   bool     `json:"isOwnerViewing"`
		ShortDescription string   `json:"shortDescription"`
		IsCrawlable      bool     `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		AverageRating     float64 `json:"averageRating"`
		AllowRatings      bool    `json:"allowRatings"`
		ViewCount         string  `json:"viewCount"`
		Author            string  `json:"author"`
		IsPrivate         bool    `json:"isPrivate"`
		IsUnpluggedCorpus bool    `json:"isUnpluggedCorpus"`
		IsLiveContent     bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
	Annotations []struct {
		PlayerAnnotationsExpandedRenderer struct {
			FeaturedChannel struct {
				StartTimeMs string `json:"startTimeMs"`
				EndTimeMs   string `json:"endTimeMs"`
				Watermark   struct {
					Thumbnails []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"thumbnails"`
				} `json:"watermark"`
				TrackingParams     string `json:"trackingParams"`
				NavigationEndpoint struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL         string `json:"url"`
							WebPageType string `json:"webPageType"`
							RootVe      int    `json:"rootVe"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					BrowseEndpoint struct {
						BrowseID string `json:"browseId"`
					} `json:"browseEndpoint"`
				} `json:"navigationEndpoint"`
				ChannelName     string `json:"channelName"`
				SubscribeButton struct {
					SubscribeButtonRenderer struct {
						ButtonText struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"buttonText"`
						Subscribed           bool   `json:"subscribed"`
						Enabled              bool   `json:"enabled"`
						Type                 string `json:"type"`
						ChannelID            string `json:"channelId"`
						ShowPreferences      bool   `json:"showPreferences"`
						SubscribedButtonText struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"subscribedButtonText"`
						UnsubscribedButtonText struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"unsubscribedButtonText"`
						TrackingParams        string `json:"trackingParams"`
						UnsubscribeButtonText struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"unsubscribeButtonText"`
						ServiceEndpoints []struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
							CommandMetadata     struct {
								WebCommandMetadata struct {
									URL      string `json:"url"`
									SendPost bool   `json:"sendPost"`
									APIURL   string `json:"apiUrl"`
								} `json:"webCommandMetadata"`
							} `json:"commandMetadata"`
							SubscribeEndpoint struct {
								ChannelIds []string `json:"channelIds"`
								Params     string   `json:"params"`
							} `json:"subscribeEndpoint,omitempty"`
							SignalServiceEndpoint struct {
								Signal  string `json:"signal"`
								Actions []struct {
									ClickTrackingParams string `json:"clickTrackingParams"`
									OpenPopupAction     struct {
										Popup struct {
											ConfirmDialogRenderer struct {
												TrackingParams string `json:"trackingParams"`
												DialogMessages []struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"dialogMessages"`
												ConfirmButton struct {
													ButtonRenderer struct {
														Style      string `json:"style"`
														Size       string `json:"size"`
														IsDisabled bool   `json:"isDisabled"`
														Text       struct {
															Runs []struct {
																Text string `json:"text"`
															} `json:"runs"`
														} `json:"text"`
														ServiceEndpoint struct {
															ClickTrackingParams string `json:"clickTrackingParams"`
															CommandMetadata     struct {
																WebCommandMetadata struct {
																	URL      string `json:"url"`
																	SendPost bool   `json:"sendPost"`
																	APIURL   string `json:"apiUrl"`
																} `json:"webCommandMetadata"`
															} `json:"commandMetadata"`
															UnsubscribeEndpoint struct {
																ChannelIds []string `json:"channelIds"`
																Params     string   `json:"params"`
															} `json:"unsubscribeEndpoint"`
														} `json:"serviceEndpoint"`
														Accessibility struct {
															Label string `json:"label"`
														} `json:"accessibility"`
														TrackingParams string `json:"trackingParams"`
													} `json:"buttonRenderer"`
												} `json:"confirmButton"`
												CancelButton struct {
													ButtonRenderer struct {
														Style      string `json:"style"`
														Size       string `json:"size"`
														IsDisabled bool   `json:"isDisabled"`
														Text       struct {
															Runs []struct {
																Text string `json:"text"`
															} `json:"runs"`
														} `json:"text"`
														Accessibility struct {
															Label string `json:"label"`
														} `json:"accessibility"`
														TrackingParams string `json:"trackingParams"`
													} `json:"buttonRenderer"`
												} `json:"cancelButton"`
												PrimaryIsCancel bool `json:"primaryIsCancel"`
											} `json:"confirmDialogRenderer"`
										} `json:"popup"`
										PopupType string `json:"popupType"`
									} `json:"openPopupAction"`
								} `json:"actions"`
							} `json:"signalServiceEndpoint,omitempty"`
						} `json:"serviceEndpoints"`
						SubscribeAccessibility struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"subscribeAccessibility"`
						UnsubscribeAccessibility struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"unsubscribeAccessibility"`
						SignInEndpoint struct {
							ClickTrackingParams string `json:"clickTrackingParams"`
						} `json:"signInEndpoint"`
					} `json:"subscribeButtonRenderer"`
				} `json:"subscribeButton"`
			} `json:"featuredChannel"`
			AllowSwipeDismiss bool   `json:"allowSwipeDismiss"`
			AnnotationID      string `json:"annotationId"`
		} `json:"playerAnnotationsExpandedRenderer"`
	} `json:"annotations"`
	PlayerConfig struct {
		AudioConfig struct {
			LoudnessDb              float64 `json:"loudnessDb"`
			PerceptualLoudnessDb    float64 `json:"perceptualLoudnessDb"`
			EnablePerFormatLoudness bool    `json:"enablePerFormatLoudness"`
		} `json:"audioConfig"`
		StreamSelectionConfig struct {
			MaxBitrate string `json:"maxBitrate"`
		} `json:"streamSelectionConfig"`
		DaiConfig struct {
			EnableServerStitchedDai bool `json:"enableServerStitchedDai"`
		} `json:"daiConfig"`
		MediaCommonConfig struct {
			DynamicReadaheadConfig struct {
				MaxReadAheadMediaTimeMs int `json:"maxReadAheadMediaTimeMs"`
				MinReadAheadMediaTimeMs int `json:"minReadAheadMediaTimeMs"`
				ReadAheadGrowthRateMs   int `json:"readAheadGrowthRateMs"`
			} `json:"dynamicReadaheadConfig"`
		} `json:"mediaCommonConfig"`
		WebPlayerConfig struct {
			WebPlayerActionsPorting struct {
				GetSharePanelCommand struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL      string `json:"url"`
							SendPost bool   `json:"sendPost"`
							APIURL   string `json:"apiUrl"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					WebPlayerShareEntityServiceEndpoint struct {
						SerializedShareEntity string `json:"serializedShareEntity"`
					} `json:"webPlayerShareEntityServiceEndpoint"`
				} `json:"getSharePanelCommand"`
				SubscribeCommand struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL      string `json:"url"`
							SendPost bool   `json:"sendPost"`
							APIURL   string `json:"apiUrl"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					SubscribeEndpoint struct {
						ChannelIds []string `json:"channelIds"`
						Params     string   `json:"params"`
					} `json:"subscribeEndpoint"`
				} `json:"subscribeCommand"`
				UnsubscribeCommand struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL      string `json:"url"`
							SendPost bool   `json:"sendPost"`
							APIURL   string `json:"apiUrl"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					UnsubscribeEndpoint struct {
						ChannelIds []string `json:"channelIds"`
						Params     string   `json:"params"`
					} `json:"unsubscribeEndpoint"`
				} `json:"unsubscribeCommand"`
				AddToWatchLaterCommand struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL      string `json:"url"`
							SendPost bool   `json:"sendPost"`
							APIURL   string `json:"apiUrl"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					PlaylistEditEndpoint struct {
						PlaylistID string `json:"playlistId"`
						Actions    []struct {
							AddedVideoID string `json:"addedVideoId"`
							Action       string `json:"action"`
						} `json:"actions"`
					} `json:"playlistEditEndpoint"`
				} `json:"addToWatchLaterCommand"`
				RemoveFromWatchLaterCommand struct {
					ClickTrackingParams string `json:"clickTrackingParams"`
					CommandMetadata     struct {
						WebCommandMetadata struct {
							URL      string `json:"url"`
							SendPost bool   `json:"sendPost"`
							APIURL   string `json:"apiUrl"`
						} `json:"webCommandMetadata"`
					} `json:"commandMetadata"`
					PlaylistEditEndpoint struct {
						PlaylistID string `json:"playlistId"`
						Actions    []struct {
							Action         string `json:"action"`
							RemovedVideoID string `json:"removedVideoId"`
						} `json:"actions"`
					} `json:"playlistEditEndpoint"`
				} `json:"removeFromWatchLaterCommand"`
			} `json:"webPlayerActionsPorting"`
		} `json:"webPlayerConfig"`
	} `json:"playerConfig"`
	Storyboards struct {
		PlayerStoryboardSpecRenderer struct {
			Spec string `json:"spec"`
		} `json:"playerStoryboardSpecRenderer"`
	} `json:"storyboards"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			Embed struct {
				IframeURL      string `json:"iframeUrl"`
				FlashURL       string `json:"flashUrl"`
				Width          int    `json:"width"`
				Height         int    `json:"height"`
				FlashSecureURL string `json:"flashSecureUrl"`
			} `json:"embed"`
			Title struct {
				SimpleText string `json:"simpleText"`
			} `json:"title"`
			Description struct {
				SimpleText string `json:"simpleText"`
			} `json:"description"`
			LengthSeconds      string   `json:"lengthSeconds"`
			OwnerProfileURL    string   `json:"ownerProfileUrl"`
			ExternalChannelID  string   `json:"externalChannelId"`
			AvailableCountries []string `json:"availableCountries"`
			IsUnlisted         bool     `json:"isUnlisted"`
			HasYpcMetadata     bool     `json:"hasYpcMetadata"`
			ViewCount          string   `json:"viewCount"`
			Category           string   `json:"category"`
			PublishDate        string   `json:"publishDate"`
			OwnerChannelName   string   `json:"ownerChannelName"`
			UploadDate         string   `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
	TrackingParams string `json:"trackingParams"`
	Attestation    struct {
		PlayerAttestationRenderer struct {
			Challenge    string `json:"challenge"`
			BotguardData struct {
				Program        string `json:"program"`
				InterpreterURL string `json:"interpreterUrl"`
			} `json:"botguardData"`
		} `json:"playerAttestationRenderer"`
	} `json:"attestation"`
	Messages []struct {
		MealbarPromoRenderer struct {
			MessageTexts []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTexts"`
			ActionButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL      string `json:"url"`
								SendPost bool   `json:"sendPost"`
								APIURL   string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						FeedbackEndpoint struct {
							FeedbackToken string `json:"feedbackToken"`
							UIActions     struct {
								HideEnclosingContainer bool `json:"hideEnclosingContainer"`
							} `json:"uiActions"`
						} `json:"feedbackEndpoint"`
					} `json:"serviceEndpoint"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						URLEndpoint struct {
							URL    string `json:"url"`
							Target string `json:"target"`
						} `json:"urlEndpoint"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"actionButton"`
			DismissButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL      string `json:"url"`
								SendPost bool   `json:"sendPost"`
								APIURL   string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						FeedbackEndpoint struct {
							FeedbackToken string `json:"feedbackToken"`
							UIActions     struct {
								HideEnclosingContainer bool `json:"hideEnclosingContainer"`
							} `json:"uiActions"`
						} `json:"feedbackEndpoint"`
					} `json:"serviceEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
			TriggerCondition    string `json:"triggerCondition"`
			Style               string `json:"style"`
			TrackingParams      string `json:"trackingParams"`
			ImpressionEndpoints []struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				CommandMetadata     struct {
					WebCommandMetadata struct {
						URL      string `json:"url"`
						SendPost bool   `json:"sendPost"`
						APIURL   string `json:"apiUrl"`
					} `json:"webCommandMetadata"`
				} `json:"commandMetadata"`
				FeedbackEndpoint struct {
					FeedbackToken string `json:"feedbackToken"`
					UIActions     struct {
						HideEnclosingContainer bool `json:"hideEnclosingContainer"`
					} `json:"uiActions"`
				} `json:"feedbackEndpoint"`
			} `json:"impressionEndpoints"`
			IsVisible    bool `json:"isVisible"`
			MessageTitle struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTitle"`
		} `json:"mealbarPromoRenderer"`
	} `json:"messages"`
	Endscreen struct {
		EndscreenRenderer struct {
			Elements []struct {
				EndscreenElementRenderer struct {
					Style string `json:"style"`
					Image struct {
						Thumbnails []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"thumbnails"`
					} `json:"image"`
					VideoDuration struct {
						SimpleText string `json:"simpleText"`
					} `json:"videoDuration"`
					Left        float64 `json:"left"`
					Width       float64 `json:"width"`
					Top         float64 `json:"top"`
					AspectRatio float64 `json:"aspectRatio"`
					StartMs     string  `json:"startMs"`
					EndMs       string  `json:"endMs"`
					Title       struct {
						Accessibility struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"accessibility"`
						SimpleText string `json:"simpleText"`
					} `json:"title"`
					Metadata struct {
						SimpleText string `json:"simpleText"`
					} `json:"metadata"`
					Endpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID string `json:"videoId"`
						} `json:"watchEndpoint"`
					} `json:"endpoint"`
					TrackingParams string `json:"trackingParams"`
					ID             string `json:"id"`
				} `json:"endscreenElementRenderer"`
			} `json:"elements"`
			StartMs        string `json:"startMs"`
			TrackingParams string `json:"trackingParams"`
		} `json:"endscreenRenderer"`
	} `json:"endscreen"`
}
