
			cmd_delete := exec.Command("sh", "-c", imgVER)
			cmd_delete.Stdout = os.Stdout
			if err := cmd_delete.Run(); err != nil {
				panic(err)
			}
			wg.Done()
		}
	default:
		fmt.Println("Check your Command")
		deletePods()
	}
	wg.Wait()
}